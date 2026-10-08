// Package handler ...
package handler

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"slices"
	"strings"
	"syscall"
	"time"

	"my-redis/app/commands"
	"my-redis/app/parser"
	"my-redis/app/replication"
)

const BufferSize = 512

type handler struct {
	fromMaster bool
	logger     *slog.Logger
	env        *commands.Env
}

func (h *handler) HandleConnection(ctx context.Context, conn net.Conn, br *bufio.Reader) {
	h.logger = h.logger.With("addr", conn.RemoteAddr())
	if br == nil {
		br = bufio.NewReader(conn)
	}
	r := parser.NewReader(br) // one buffer shared by parser and watcher
	sess := &commands.Session{
		Queue:     make([]commands.QueueItem, 0),
		Role:      h.env.Repl.Role(),
		WatchList: make([]commands.WatchItem, 0),
	}
	sess.Watch = watchDisconnect(conn, br)
	if sess.Role == replication.MasterRole {
		sess.Register = func() {
			// h.env.Repl.Register(conn, replication.EmptyDB, r)
			sess.Detached = true
		}
	}

	stop := context.AfterFunc(ctx, func() {
		if err := conn.Close(); err != nil {
			h.logger.Error("error closing connection", "err", err)
		}
		h.logger.Debug("closing handler conn after ctx")
	}) // parent ctx, no goroutine
	defer stop()
	defer func() {
		if !sess.Detached { // a detached conn belongs to replication
			h.logger.Debug("closing connection")
			if err := conn.Close(); err != nil {
				h.logger.Error("closing conn", "err", err)
			}
		}
	}()

	for {
		arr, err := r.ReadArrays()
		if err != nil {
			h.logReadError(err)
			return
		}
		if len(arr) == 0 {
			h.logger.Warn("got empty commands array from client without error")
		}
		for _, args := range arr {
			res := h.Dispatch(ctx, args, sess)
			if sess.Detached {
				// run() has returned, so the Peek goroutine is gone and the deadline is reset.
				h.env.Repl.Register(conn, replication.EmptyDB, r)
				return
			}
			if h.isReplica() && h.fromMaster {
				h.AddProcessed(args)
			}
			if sess.Detached { // check BEFORE writing
				h.logger.Debug("replica conn, move to detached mode", "cmd", args[0])
				return
			}
			if res == nil {
				continue
			}
			if h.fromMaster && !isGetAck(args) {
				h.logger.Debug("from master, don't response", "is_get_ack", isGetAck(args), "args", args)
				continue
			}
			response := parser.Encode(res)
			if len(response) == 0 {
				h.logger.Debug("attempt to send an empty response", "args", args, "res", res)
			}
			if _, err := conn.Write(response); err != nil {
				h.logger.Debug("write failed", "addr", conn.RemoteAddr(), "err", err)
				return
			}
		}
	}
}

func isGetAck(args []string) bool {
	res := slices.Equal(args, []string{"REPLCONF", "GETACK", "*"})
	// slog.Debug("ack?", "args", args, "res", res)
	return res
}

func (h *handler) AddProcessed(args []string) {
	cmd := parser.CommandFromStrings(args...)
	bytes, _ := cmd.Serialize()
	h.env.Repl.AddProcessed(len(bytes))
	// h.logger.Debug("processed bytes", "n", len(bytes), "args", args)
}

func (h *handler) Dispatch(ctx context.Context, args []string, sess *commands.Session) parser.Serializable {
	h.logger.Info("dispatching", "args", args)
	name := strings.ToUpper(args[0])
	args = args[1:]
	switch name {
	case "MULTI":
		if sess.InMulti {
			sess.WatchList = sess.WatchList[:0]
			return parser.SimpleError("ERR MULTI inside MULTI is not allowed")
		}
		sess.InMulti = true

		return parser.SimpleString("OK")

	case "WATCH":
		if sess.InMulti {
			return parser.SimpleError("ERR WATCH inside MULTI is not allowed")
		}
	case "DISCARD":
		if !sess.InMulti {
			return parser.SimpleError("ERR DISCARD without MULTI")
		}
		sess.FinishTx()

		return parser.SimpleString("OK")
	case "EXEC":
		if !sess.InMulti {
			return parser.SimpleError("ERR EXEC without MULTI")
		}

		return h.exec(ctx, sess)
	}
	cmdSpec, ok := h.env.Reg.Lookup(name)
	if !ok {
		return parser.SimpleError(fmt.Sprintf("command not found: %v", name))
	}

	command, err := cmdSpec.ParseArgs(args)
	if err != nil {
		h.logger.Error("Error parsing command", "error", err.Error())
		return parser.SimpleError(err.Error())
	}
	if sess.InMulti {
		sess.Queue = append(sess.Queue, commands.QueueItem{Cmd: command, Spec: cmdSpec, Args: args})
		return parser.SimpleString("QUEUED")
	}

	return h.run(ctx, command, &cmdSpec, args, sess)
}

func (h *handler) exec(ctx context.Context, sess *commands.Session) parser.Serializable {
	h.env.Wg.Add(1)
	defer h.env.Wg.Done()
	defer sess.FinishTx()
	for _, w := range sess.WatchList {
		v := h.env.Coordinator.GetVersions(w.Key)[0]
		if v != w.Version {
			return parser.NullArray(0)
		}
	}

	res := make(parser.Array[parser.Serializable], len(sess.Queue))
	for i, cmd := range sess.Queue {
		r := h.run(ctx, cmd.Cmd, &cmd.Spec, cmd.Args, sess)
		res[i] = r
	}
	return res
}

func (h *handler) run(ctx context.Context, cmd commands.Command, spec *commands.Spec, args []string, sess *commands.Session) parser.Serializable {
	if sess.Watch != nil {
		var stop func()
		ctx, stop = sess.Watch(ctx)
		defer stop()
	}
	if h.isReplica() && !h.fromMaster && (spec.Flags&commands.FlagWrite != 0) {
		return parser.SimpleError("ERR cannot write to replica")
	}
	r, err := cmd.Execute(ctx, h.env, sess)
	if err != nil {
		return parser.SimpleError(err.Error())
	}
	if h.isMaster() && h.env.Cfg.AppendOnly == "yes" && (spec.Flags&commands.FlagWrite != 0) {
		cmdArgs := slices.Insert(args, 0, spec.Name)
		toSend := parser.CommandFromStrings(cmdArgs...)
		bytes, _ := toSend.Serialize()
		h.env.ManifestChan <- string(bytes)
	}
	if h.isMaster() && (spec.Flags.Propagate() || isGetAck(args)) {
		cmdArgs := slices.Insert(args, 0, spec.Name)
		toSend := parser.CommandFromStrings(cmdArgs...)
		bytes, err := toSend.Serialize()
		if err != nil {
			h.logger.Warn("try to propagate zero bytes", "args", cmdArgs, "bytes", string(bytes))
		}
		h.env.Repl.Propagate(bytes)
		h.logger.Debug("propagated", "args", cmdArgs)
	}
	return r
}

func (h *handler) isReplica() bool {
	return h.env.Repl.Role() == replication.ReplicaRole
}

func (h *handler) isMaster() bool {
	return h.env.Repl.Role() == replication.MasterRole
}

func watchDisconnect(conn net.Conn, br *bufio.Reader) func(context.Context) (context.Context, func()) {
	return func(ctx context.Context) (context.Context, func()) {
		ctx, cancel := context.WithCancel(ctx)
		done := make(chan struct{})
		go func() {
			defer close(done)
			if _, err := br.Peek(1); err != nil && !errors.Is(err, os.ErrDeadlineExceeded) {
				cancel() // EOF or reset: the client is gone
			}
		}()
		return ctx, func() {
			_ = conn.SetReadDeadline(time.Now()) // unblock Peek
			<-done
			_ = conn.SetReadDeadline(time.Time{})
			cancel()
		}
	}
}

func (h *handler) logReadError(err error) {
	switch {
	case errors.Is(err, io.EOF), errors.Is(err, net.ErrClosed):
		// clean disconnect or our own Close: nothing to log
	case errors.Is(err, io.ErrUnexpectedEOF), errors.Is(err, syscall.ECONNRESET):
		h.logger.Debug("client dropped mid-command", "err", err)
	default:
		h.logger.Warn("input parsing failed", "err", err)
	}
}
