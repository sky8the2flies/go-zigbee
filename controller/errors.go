package controller

import "errors"

var (
	ErrEngineClosed = errors.New("engine closed")
	ErrAdapterStart = errors.New("adapter start failed")
)
