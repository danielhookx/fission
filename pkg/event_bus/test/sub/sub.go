package main

import (
	"fmt"

	eventbus "github.com/danielhookx/fission/pkg/event_bus"
)

func SayHello(name string) {
	fmt.Printf("hello %s\n", name)
}

func main() {
	rawURL := "tcp://:7634"
	remoteURL := "tcp://localhost:7633"

	bus := eventbus.NewEventBus(eventbus.WithRemoteMode(eventbus.CreateNewIPCHandler(rawURL, remoteURL)))
	bus.Subscribe("test-hello", SayHello)

	select {}
}
