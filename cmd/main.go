package main

import (
	"log/slog"
	"os"
	"time"

	"go-zigbee/internal/controller"
	"go-zigbee/internal/parser"
	"go-zigbee/internal/reader"
	"go-zigbee/internal/transport"
	"go-zigbee/pkg/pubsub"

	"github.com/golang-cz/devslog"
	"github.com/tarm/serial"
)

func main() {
	logger := slog.New(devslog.NewHandler(os.Stdout, &devslog.Options{
		SortKeys:          true,
		NewLineAfterLog:   true,
		StringerFormatter: true,
		HandlerOptions: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}))
	slog.SetDefault(logger)

	transport := transport.NewSerialPort(&serial.Config{
		Name:        "/dev/ttyAMA0",
		Baud:        38400,
		ReadTimeout: time.Second,
	})

	agent := pubsub.NewAgent[any]()
	parser := parser.NewParser()
	reader := reader.NewReader(parser, agent)

	controller := controller.NewController(transport, reader)
	err := controller.Start()
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	ca := agent.Subscribe(pubsub.Unknown)

	defer agent.Close()
	for chunk := range ca {
		slog.Debug("chunk", slog.Any("chunk", chunk))
	}
}
