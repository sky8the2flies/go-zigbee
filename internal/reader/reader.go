package reader

import (
	"log/slog"

	"go-zigbee/internal/controller"
	"go-zigbee/pkg/pubsub"
)

type Reader struct {
	Parser      Parser
	Agent       *pubsub.Agent[any]
	errorSignal chan error
	chunkSignal chan []byte
}

func NewReader(parser Parser, agent *pubsub.Agent[any]) *Reader {
	return &Reader{Parser: parser, Agent: agent}
}

func (c *Reader) Start(port controller.Port) {
	defer port.Close()

	c.errorSignal = make(chan error, 1)
	c.chunkSignal = make(chan []byte, 1024)

	inputBuf := make([]byte, 1024)

	for {
		n, err := port.Read(inputBuf)
		if err != nil && n > 0 {
			// error occurred while data was read
			c.errorSignal <- err
			slog.Error("Error reading data", slog.String("error", err.Error()))
			continue
		}

		if n > 0 {
			d, chunk, err := c.Parser.Chunk(inputBuf[:n])
			if err != nil {
				c.errorSignal <- err
				slog.Error("Error transforming data", slog.String("error", err.Error()))
				continue
			}

			if !d {
				continue
			}

			err = c.Parser.Frame(chunk)
			if err != nil {
				c.errorSignal <- err
				slog.Error("Error framing data", slog.String("error", err.Error()), slog.Any("chunk", chunk))
				continue
			}

			c.chunkSignal <- chunk
		}
	}
}
