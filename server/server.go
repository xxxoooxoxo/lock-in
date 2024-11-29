package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/coder/websocket"
	"golang.org/x/time/rate"
)

type echoServer struct {
	logf            func(f string, v ...interface{})
	connectionCount atomic.Int64
}

func (s *echoServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.connectionCount.Add(1)
	s.logf("New connection attempt from %v (Total connections: %d)", r.RemoteAddr, s.connectionCount.Load())

	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols: []string{"echo"},
	})
	if err != nil {
		s.logf("Connection failed from %v (Total connections: %d)", r.RemoteAddr, s.connectionCount.Load())
		s.logf("%v", err)
		return
	}
	defer func() {
		c.CloseNow()
		s.connectionCount.Add(-1)
		s.logf("Connection closed with %v (Total connections: %d)", r.RemoteAddr, s.connectionCount.Load())
	}()

	if c.Subprotocol() != "echo" {
		s.logf("Client %v disconnected: invalid subprotocol", r.RemoteAddr)
		s.logf("Connection closed due to invalid subprotocol (Total connections: %d)", s.connectionCount.Load())
		c.Close(websocket.StatusPolicyViolation, "client must speak the echo subprotocol")
		return
	}

	l := rate.NewLimiter(rate.Every(time.Millisecond*100), 10)
	for {
		err = echo(r.Context(), c, l)
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			s.logf("Normal closure from %v (Total connections: %d)", r.RemoteAddr, s.connectionCount.Load())
			return
		}
		if err != nil {
			s.logf("failed to echo with %v: %v", r.RemoteAddr, err)
			s.logf("Connection closed due to error (Total connections: %d)", s.connectionCount.Load())
			return
		}
	}
}

func (s *echoServer) GetConnectionCount() int64 {
	return s.connectionCount.Load()
}

func echo(ctx context.Context, c *websocket.Conn, l *rate.Limiter) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	err := l.Wait(ctx)
	if err != nil {
		return err
	}

	typ, r, err := c.Reader(ctx)
	if err != nil {
		return err
	}

	w, err := c.Writer(ctx, typ)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, r)
	if err != nil {
		return fmt.Errorf("failed to io.Copy: %w", err)
	}

	err = w.Close()
	return err
}
