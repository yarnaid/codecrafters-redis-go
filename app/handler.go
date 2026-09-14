package main

import (
	"errors"
	"io"

	"fmt"

	"github.com/codecrafters-io/redis-starter-go/app/commands"
	"github.com/codecrafters-io/redis-starter-go/app/parser"
)

const BUFFER_SIZE = 512

type RW interface {
	io.Reader
	io.Writer
	Close() error
}

type Handler struct {
	conn RW
}

func (h *Handler) SetConn(conn RW) {
	h.conn = conn
}

func (h *Handler) handle_connection() {
	logger.Info("Handling new connection...")
	defer func() {
		h.conn.Close()
		logger.Debug("Connection closed")
	}()

	for {
		var buf [BUFFER_SIZE]byte
		n, err := h.conn.Read(buf[:])
		if err != nil {
			logger.Error("Error reading from connection", "error", err.Error())
			break
		}
		logger.Debug("data received", "data", string(buf[:n]))

		h.process_command(buf[:n])

	}
}

func (h *Handler) process_command(buf []byte) {
	input_data, err := parser.Deserialize(buf)
	if err != nil {
		logger.Error("Error parsing input data", "error", err.Error())
		h.send_error_response("Error parsing input data", err)
		return
	}

	command_arr, err := stringify(input_data.([]parser.Serializable))
	if err != nil {
		msg := "Cannot convert input to strings array"
		logger.Error(msg, "data", input_data, "type", fmt.Sprintf("%#v\n", input_data))
		h.send_error_response(msg, err)
		return
	}
	command, err := commands.FromInput(command_arr)
	if err != nil {
		logger.Error("Error parsing command", "error", err.Error())
		h.send_error_response("Error parsing command", err)
		return
	}

	response, err := command.Execute()
	if err != nil {
		logger.Error("Error executing command", "error", err.Error())
		h.send_error_response("Error executing command", err)
		return
	}
	h.send_response(response)
}

func (h *Handler) send_error_response(message string, err error) {
	msg := "Error: " + message
	if err != nil {
		msg += " - " + err.Error()
	}
	logger.Error("Sending error response", "message", message, "error", err.Error())
	h.send_response(parser.SimpleError(msg))
}

func (h *Handler) send_response(response parser.Serializable) {
	data, err := response.Serialize()
	if err != nil {
		h.send_error_response("Error serializing response", err)
		return
	}

	_, err = h.conn.Write(data)
	if err != nil {
		logger.Error("Error writing response", "error", err.Error())
	}
}

func stringify(input []parser.Serializable) ([]string, error) {
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
	logger.Debug("Stringify", "input", input, "output", args)
	return args, nil
}
