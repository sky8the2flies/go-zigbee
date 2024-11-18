package transport

import (
	"testing"

	"github.com/tarm/serial"
)

func TestSerialTransportOpen(t *testing.T) {
	mockPort := "/dev/pts/1"
	config := &serial.Config{Name: mockPort, Baud: 9600}
	transport := NewSerialPort(config)

	// Test closing when port is not open
	err := transport.Close()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Test closing when port is open
	err = transport.Open()
	if err != nil {
		t.Fatalf("failed to open port: %v", err)
	}

	err = transport.Close()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if transport.open {
		t.Errorf("expected port to be closed, but it is open")
	}
}
