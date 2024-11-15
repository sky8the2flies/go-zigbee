package transport

import "github.com/tarm/serial"

type SerialTransport struct {
	Port *serial.Port
}

func NewSerialTransport(config *serial.Config) (*SerialTransport, error) {
	port, err := serial.OpenPort(config)
	if err != nil {
		return nil, err
	}

	return &SerialTransport{
		Port: port,
	}, nil
}

func (t *SerialTransport) Read(p []byte) (n int, err error) {
	return t.Port.Read(p)
}

func (t *SerialTransport) Write(p []byte) (n int, err error) {
	return t.Port.Write(p)
}

func (t *SerialTransport) Close() error {
	return t.Port.Close()
}
