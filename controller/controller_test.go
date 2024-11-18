package controller

import (
	"testing"
	"time"
)

type MockAdapter struct {
	startErr error
	readErr  error
}

func (m *MockAdapter) Start() error {
	return m.startErr
}

func (m *MockAdapter) Read() error {
	return m.readErr
}

func (m *MockAdapter) Close() error {
	return nil
}

func TestEngine_Start(t *testing.T) {
	adapter := &MockAdapter{}
	controller := NewController(adapter)

	go func() {
		time.Sleep(1 * time.Millisecond)
		controller.Stop()
	}()

	err := controller.Start()
	if err != ErrEngineClosed {
		t.Errorf("expected %v, got %v", ErrEngineClosed, err)
	}
}
