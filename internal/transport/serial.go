package transport

import "github.com/tarm/serial"

type SerialTransport struct {
	Config *serial.Config
	port   *serial.Port
	open   bool
}

func NewSerialPort(config *serial.Config) *SerialTransport {
	return &SerialTransport{
		Config: config,
		open:   false,
	}
}

func (t *SerialTransport) Open() error {
	port, err := serial.OpenPort(t.Config)
	if err != nil {
		return err
	}
	t.port = port
	t.open = true
	return nil
}

func (t *SerialTransport) Read(p []byte) (n int, err error) {
	if !t.open {
		return 0, ErrPortNotOpen
	}
	return t.port.Read(p)
}

func (t *SerialTransport) Write(p []byte) (n int, err error) {
	if !t.open {
		return 0, ErrPortNotOpen
	}
	return t.port.Write(p)
}

func (t *SerialTransport) Close() error {
	err := t.port.Close()
	if err != nil {
		return err
	}

	t.open = false
	return nil
}
