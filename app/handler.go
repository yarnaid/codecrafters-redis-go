package main

import (
	"net"
)

type Handler struct {
	conn net.Conn
}

func (h *Handler) handle_connection() {
	logger.Info("Handling new connection...")
	defer func() {
		h.conn.Close()
		logger.Debug("Connection closed")
	}()

	for {
		var buf [512]byte
		_, err := h.conn.Read(buf[:])
		if err != nil {
			logger.Error("Error reading from connection", "error", err.Error())
			break
		}

		_, err = h.conn.Write([]byte(fmt_response("PONG")))
		if err != nil {
			logger.Error("Error writing response", "error", err.Error())
		}

	}
}
