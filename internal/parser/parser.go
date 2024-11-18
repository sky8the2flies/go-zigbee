package parser

import (
	"go-zigbee-herdsman/internal/decoder"
	"go-zigbee-herdsman/pkg/buffer"
	"log"
	"log/slog"
)

type Parser struct {
	decoder *decoder.Decoder
}

func NewParser() *Parser {
	decoder := decoder.NewDecoder()
	return &Parser{decoder: decoder}
}

func (p *Parser) Chunk(inputBuf []byte) (bool, []byte, error) {
	d, chunk := p.decoder.Decode(inputBuf)
	return d, chunk, nil
}

func (p *Parser) Frame(chunk []byte) error {
	type data struct {
		CommandID   uint8
		Sequence    uint8
		Status      uint8
		FrameLength uint16
		FrameBuffer *buffer.Scannable
	}
	temp := map[int]func(data) error{
		-1: func(data data) error {
			log.Printf("Unknown command ID %d %#x %b", data.CommandID, data.CommandID, data.CommandID)
			log.Printf("buffer %v", data.FrameBuffer.Raw())
			return nil
		},
		0x0e: func(data data) error {
			log.Printf("framebuffer %v", data.FrameBuffer.Raw())
			return nil
		},
	}
	scanbuf := buffer.NewScannable(chunk)

	commandID, err := scanbuf.GetUint8(0, true)
	if err != nil {
		slog.Error("Error reading command ID", slog.String("error", err.Error()))
		return err
	}
	seq, err := scanbuf.GetUint8(1, true)
	if err != nil {
		slog.Error("Error reading sequence number", slog.String("error", err.Error()))
		return err
	}
	status, err := scanbuf.GetUint8(2, true)
	if err != nil {
		slog.Error("Error reading status", slog.String("error", err.Error()))
		return err
	}
	frameLength, err := scanbuf.GetUint16(3, true)
	if err != nil {
		slog.Error("Error reading frame length", slog.String("error", err.Error()))
		return err
	}
	info := data{CommandID: commandID, Sequence: seq, Status: status, FrameLength: frameLength, FrameBuffer: scanbuf}
	if _, ok := temp[int(commandID)]; !ok {
		temp[-1](info)
		return nil
	}
	temp[int(commandID)](info)
	return nil
}
