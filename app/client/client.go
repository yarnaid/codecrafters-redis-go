// Package client ...
package client

import (
	"bufio"
	"net"
	"sync"
	"time"
)

const writeTimeout = 5 * time.Second

type Client struct {
	conn      net.Conn
	Addr      net.Addr
	out       chan []byte
	done      chan struct{}
	closeOnce sync.Once
}

func (c *Client) Close() {
	c.closeOnce.Do(func() {
		close(c.done)
		_ = c.conn.Close()
	})
}

func (c *Client) Done() <-chan struct{} { return c.done }

func (c *Client) writeLoop() {
	w := bufio.NewWriter(c.conn)
	for {
		select {
		case frame := <-c.out:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			if _, err := w.Write(frame); err != nil {
				c.Close()
				return
			}
			// Coalesce: flush only when the queue is drained.
			if len(c.out) == 0 {
				if err := w.Flush(); err != nil {
					c.Close()
					return
				}
			}
		case <-c.done:
			return
		}
	}
}

func (c *Client) Reply(frame []byte) bool {
	select {
	case c.out <- frame:
		return true
	case <-c.done:
		return false
	}
}

func NewClient(conn net.Conn) *Client {
	c := &Client{
		conn: conn,
		out:  make(chan []byte),
		done: make(chan struct{}),
		Addr: conn.RemoteAddr(),
	}
	go c.writeLoop()
	return c
}
