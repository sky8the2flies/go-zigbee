package adapter_deconz

import "go-zigbee/adapter"

type Adapter struct {
	transport adapter.Transport
	parser    *Parser

	bufferCap int

	buffer []byte
}

func NewAdapter(trasnsport adapter.Transport, opts ...Option) *Adapter {
	if len(opts) == 0 {
		opts = []Option{WithBufferCap(1024)}
	}

	p := NewParser()
	a := &Adapter{transport: trasnsport, parser: p}

	for _, opt := range opts {
		opt(a)
	}

	a.buffer = make([]byte, a.bufferCap)

	return a
}

func (a *Adapter) Start() error {
	if err := a.transport.Open(); err != nil {
		return err
	}
	return nil
}

func (a *Adapter) Read() error {
	inputBuf := make([]byte, a.bufferCap)
	n, err := a.transport.Read(inputBuf)
	if err != nil && n > 0 {
		return err
	}

	if n > 0 {
		// Write current input buffer to the adapter buffer
		a.buffer = append(a.buffer, inputBuf[:n]...)
		complete, _, err := a.parser.Decode(a.buffer, a.bufferCap)
		if err != nil {
			return err
		}

		if !complete {
			return nil
		}

	}

	return nil
}

func (a *Adapter) Close() error {
	return a.transport.Close()
}
