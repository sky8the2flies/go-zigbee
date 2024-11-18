package controller

import (
	"context"
	"log/slog"
)

type Controller struct {
	adapter Adapter

	errorSignal chan error

	ctx    context.Context
	cancel context.CancelFunc
}

func NewController(adapter Adapter, opts ...Option) *Controller {
	if len(opts) == 0 {
		opts = []Option{}
	}

	c := &Controller{adapter: adapter}

	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Controller) Start() error {
	if c.adapter.Start() != nil {
		return ErrAdapterStart
	}

	ctx, cancel := context.WithCancel(context.Background())
	c.ctx = ctx
	c.cancel = cancel

	c.errorSignal = make(chan error, 1)

	return c.engine()
}

func (c *Controller) Stop() error {
	c.cancel()
	return c.adapter.Close()
}

func (c *Controller) engine() error {
	for {
		select {
		case <-c.ctx.Done():
			return ErrEngineClosed
		case err := <-c.errorSignal:
			slog.Error("Error from signal", slog.String("error", err.Error()))
		default:
			if err := c.reader(); err != nil {
				c.errorSignal <- err
			}
		}
	}
}

func (c *Controller) reader() error {
	return c.adapter.Read()
}
