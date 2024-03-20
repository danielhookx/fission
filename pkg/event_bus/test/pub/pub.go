package main

import (
	"time"

	eventbus "github.com/danielhookx/fission/pkg/event_bus"
)

func main() {
	rawURL := "tcp://:7633"
	remoteURL := "tcp://localhost:7634"

	bus := eventbus.NewEventBus(eventbus.WithRemoteMode(eventbus.CreateNewIPCHandler(rawURL, remoteURL)))

	time.Sleep(time.Second * 10)
	ticker := time.NewTicker(time.Second * 2)
	for {
		select {
		case <-ticker.C:
			bus.Publish("test-hello", "daniel")
		}
	}
}
