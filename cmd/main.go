package main

import (
	"go-zigbee-herdsman/internal/controller"
	"go-zigbee-herdsman/internal/transport"
	"time"

	"github.com/tarm/serial"
)

func main() {
	transport := transport.NewSerialPort(&serial.Config{
		Name:        "/dev/ttyAMA0",
		Baud:        38400,
		ReadTimeout: time.Second,
	})

	controller := controller.NewController(transport)
	controller.Start()
}
