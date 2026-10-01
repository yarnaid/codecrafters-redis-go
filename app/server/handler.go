package server

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"sync"

	"my-redis/app/commands"
	"my-redis/app/parser"
	"my-redis/app/storage"
)

const BUFFER_SIZE = 512

type RW interface {
	io.Reader
	io.Writer
	Close() error
}
type WatchItem struct {
	Key     string
	Version int
}
type Handler struct {
	Conn               RW
	Coordinator        *storage.Coordinator
	Queue              []commands.Command
	TransactionStarted bool
	WatchList          []WatchItem
	TxWait             *sync.WaitGroup
	Cfg                *config
	serverState        *ServerState
}

func (h *Handler) SetConn(conn RW) {
	h.Conn = conn
}

func (h *Handler) HandleConnection() {
	logger.Info("Handling new connection...")
	defer h.Conn.Close()
	h.Queue = make([]commands.Command, 0)

	for {
		var buf [BUFFER_SIZE]byte
		n, err := h.Conn.Read(buf[:])
		if err != nil {
			logger.Error("Error reading from connection", "error", err.Error())
			break
		}
		logger.Debug("[Handler]", "data", string(buf[:n]))

		h.TxWait.Wait()
		h.processCommand(buf[:n])

	}
}

func (h *Handler) processCommand(buf []byte) {
	inputData, err := parser.Deserialize(buf)
	if err != nil {
		logger.Error("Error parsing input data", "error", err.Error())
		h.sendErrorResponse("Error parsing input data", err)
		return
	}

	commandArr, err := toStringsSlice(inputData.([]parser.Serializable))
	if err != nil {
		msg := "Cannot convert input to strings array"
		logger.Error(msg, "data", inputData, "type", fmt.Sprintf("%#v\n", inputData))
		h.sendErrorResponse(msg, err)
		return
	}
	command, err := h.FromInput(commandArr)
	if err != nil {
		logger.Error("Error parsing command", "error", err.Error())
		h.sendErrorResponse("Error parsing command", err)
		return
	}

	var response parser.Serializable
	switch c := command.(type) {
	case *commands.WatchCommand:
		response, err = c.Execute()
		if err == nil {
			for i := range c.Keys {
				h.WatchList = append(h.WatchList, WatchItem{c.Keys[i], c.Versions[i]})
				slog.Debug("[handler]start watching", "item", h.WatchList[len(h.WatchList)-1])
			}
		} else {
			slog.Error("[handler] add watch keys error", "err", err, "keys", c.Keys)
			h.WatchList = h.WatchList[:0]
		}
	case *commands.MultiCommand:
		h.TransactionStarted = true
		slog.Debug("[Handler] transaction started", "tx", h.TransactionStarted)
		response, err = command.Execute()
		if err != nil {
			h.finishTx()
		}
	case *commands.UnwatchCommand:
		response, err = command.Execute()
		h.WatchList = h.WatchList[:0]
	case *commands.InfoCommand:
		response, err = h.serverState.toBulkString(), nil
	case *commands.DiscardCommand:
		response, err = command.Execute()
		h.finishTx()
	case *commands.ExecCommand:
		slog.Debug("[Handler] transaction exec", "tx", h.TransactionStarted, "cmd.tx", c.InTransaction)
		h.TxWait.Add(1)
		defer h.TxWait.Done()
		h.TransactionStarted = false
		if len(h.WatchList) > 0 {
			for _, w := range h.WatchList {
				slog.Debug("[handler] validating watch keys", "watch", h.WatchList)

				v := h.Coordinator.GetVersions(w.Key)[0]
				if v != w.Version {
					h.finishTx()
					h.sendResponse(parser.NullArray(0))
					return
				}
			}
			slog.Debug("[handler] all watch keys are valid")

		} else {
			slog.Debug("[handler] no watch keys found")
		}
		response, err = command.Execute()
		if err == nil {
			responses := make(parser.Array[parser.Serializable], len(h.Queue))
			for i, cmd := range h.Queue {
				r, e := cmd.Execute()
				if e != nil {
					responses[i] = parser.SimpleError(e.Error())
				} else {
					responses[i] = r
				}
			}
			response = responses
		}
		h.finishTx()
	default:
		if h.TransactionStarted {
			h.Queue = append(h.Queue, command)
			response = parser.SimpleString("QUEUED")
		} else {
			response, err = command.Execute()
		}
	}

	if err != nil {
		logger.Error("Error executing command", "error", err.Error())
		h.sendResponse(parser.SimpleError(err.Error()))
		return
	}
	h.sendResponse(response)
}

func (h *Handler) sendErrorResponse(message string, err error) {
	msg := "Error: " + message
	if err != nil {
		msg += " - " + err.Error()
	}
	logger.Error("Sending error response", "message", message, "error", err.Error())
	h.sendResponse(parser.SimpleError(msg))
}

func (h *Handler) sendResponse(response parser.Serializable) {
	logger.Debug("[handler] sending response", "response", response)
	data, err := response.Serialize()
	if err != nil {
		h.sendErrorResponse("Error serializing response", err)
		return
	}

	_, err = h.Conn.Write(data)
	if err != nil {
		logger.Error("Error writing response", "error", err.Error())
	}
}

func toStringsSlice(input []parser.Serializable) ([]string, error) {
	args := make([]string, len(input))
	for i, s := range input {
		switch v := s.(type) {
		case parser.BulkString:
			args[i] = string(v)
		case parser.SimpleString:
			args[i] = string(v)
		default:
			msg := "unexpected element type in command array"
			return nil, errors.New(msg)
		}
	}
	// logger.Debug("Stringify", "input", input, "output", args)
	return args, nil
}

func (h *Handler) finishTx() {
	h.TransactionStarted = false
	clear(h.Queue)
	h.Queue = h.Queue[:0]
	h.WatchList = h.WatchList[:0]
}
