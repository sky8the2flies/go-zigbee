package controller

import (
	"bytes"
	"encoding/binary"
	"log/slog"
	"time"
)

type Controller struct {
	Port   Port
	Reader Reader
}

func NewController(port Port, reader Reader) *Controller {
	slog.Debug("Creating new controller")
	return &Controller{Port: port, Reader: reader}
}

func calcCRC(data []byte) (byte, byte) {
	crc := 0
	for _, b := range data {
		crc += int(b)
	}
	crc0 := (^crc + 1) & 0xff
	crc1 := ((^crc + 1) >> 8) & 0xff
	return byte(crc0), byte(crc1)
}

func encode(data []byte) []byte {
	crc0, crc1 := calcCRC(data)
	data = append(data, crc0, crc1)
	data = append([]byte{192}, data...)
	data = append(data, 192)
	return data
}

func initBuf(seq, requestId, paramLen, paramId int) []byte {
	payloadLen := 1 + paramLen
	frameLen := 7 + payloadLen

	fl1 := byte(frameLen & 0xff)
	fl2 := byte(frameLen >> 8)

	pl1 := byte(payloadLen & 0xff)
	pl2 := byte(payloadLen >> 8)

	return []byte{byte(requestId), byte(seq), 0x00, fl1, fl2, pl1, pl2, byte(paramId)}
}

func (c *Controller) Start() error {
	slog.Debug("Opening port")
	err := c.Port.Open()
	if err != nil {
		slog.Error("Error opening port", slog.String("error", err.Error()))
		return err
	}

	slog.Debug("Starting reader from controller")
	go c.Reader.Start(c.Port)

	// readChannelMask := encode([]byte{0x0a, 0, 0, 0x08, 0x00, 0x01, 0x00, 0x1c})

	// time.Sleep(1 * time.Second)

	// writeAllowConnections := initBuf(0x0b, 1, 0x21)
	// writeAllowConnections = append(writeAllowConnections, 0x00)
	// writeAllowConnections = encode(writeAllowConnections)
	paramId := 0x0a
	paramLen := 4
	writeChannelMask := initBuf(0, 0x0b, paramLen, paramId)
	buffer := bytes.NewBuffer(writeChannelMask)
	binary.Write(buffer, binary.LittleEndian, uint32(0x800))

	writeChannelMask = encode(buffer.Bytes())

	slog.Debug("Writing data", slog.Any("data", writeChannelMask))
	_, err = c.Port.Write(writeChannelMask)
	if err != nil {
		slog.Error("Error writing data", slog.String("error", err.Error()))
	}

	time.Sleep(2 * time.Second)

	writeNetReq := encode([]byte{0x08, 1, 0, 6, 0, 0x00})

	slog.Debug("Writing data", slog.Any("data", writeNetReq))
	_, err = c.Port.Write(writeNetReq)
	if err != nil {
		slog.Error("Error writing data", slog.String("error", err.Error()))
	}

	time.Sleep(2 * time.Second)

	writeNetReq = encode([]byte{0x08, 2, 0, 6, 0, 0x02})

	slog.Debug("Writing data", slog.Any("data", writeNetReq))
	_, err = c.Port.Write(writeNetReq)
	if err != nil {
		slog.Error("Error writing data", slog.String("error", err.Error()))
	}

	time.Sleep(2 * time.Second)

	paramId = 0x21
	paramLen = 1
	writePermitJoin := initBuf(3, 0x0b, paramLen, paramId)
	writePermitJoin = append(writePermitJoin, 20)

	writePermitJoin = encode(writePermitJoin)

	slog.Debug("Writing data", slog.Any("data", writePermitJoin))

	_, err = c.Port.Write(writePermitJoin)
	return nil
}
