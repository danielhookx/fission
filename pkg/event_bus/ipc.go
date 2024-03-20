package eventbus

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/url"

	"github.com/google/uuid"
)

type IPC struct {
	id        string
	rawURL    string
	remoteURL string
	bus       Eventbus
}

func CreateNewIPCHandler(rawURL, remoteURL string) NewRemoteBusHandler {
	return func(bus Eventbus) BusSubscriber {
		return NewIPC(rawURL, remoteURL, bus)
	}
}

func NewIPC(rawURL, remoteURL string, bus Eventbus) *IPC {
	i := &IPC{
		id:        uuid.New().String(),
		rawURL:    rawURL,
		remoteURL: remoteURL,
		bus:       bus,
	}
	rpc.Register(i)
	// Registers an HTTP handler for RPC messages
	rpc.HandleHTTP()
	// Start listening for the requests on port 1234

	u, err := url.Parse(rawURL)
	if err != nil {
		log.Fatal("parse url error: ", err)
	}

	listener, err := net.Listen(u.Scheme, fmt.Sprintf(":%s", u.Port()))
	if err != nil {
		log.Fatal("Listener error: ", err)
	}
	// Serve accepts incoming HTTP connections on the listener l, creating
	// a new service goroutine for each. The service goroutines read requests
	// and then call handler to reply to them
	go http.Serve(listener, nil)
	return i
}

type SubArgs struct {
	RemoteURL string
	Topic     string
}

type SubReply struct {
}

type PubArgs struct {
	Topic string
	data  any
}

type PubReply struct {
}

func (i *IPC) Subscribe(topic string, fn interface{}) error {
	// DialHTTP connects to an HTTP RPC server at the specified network
	remote, err := url.Parse(i.remoteURL)
	if err != nil {
		log.Fatal("parse url error: ", err)
		return err
	}
	client, err := rpc.DialHTTP(remote.Scheme, remote.Host)
	if err != nil {
		log.Fatal("Client connection error: ", err)
		return err
	}

	reply := &SubReply{}
	err = client.Call("IPC.RPCSubscribe", &SubArgs{
		RemoteURL: i.rawURL,
		Topic:     topic,
	}, &reply)
	if err != nil {
		log.Fatal("Client invocation error: ", err)
		return err
	}
	return nil
}

func (i *IPC) SubscribeSync(topic string, fn interface{}) error {
	return nil
}

func (i *IPC) Unsubscribe(topic string, handler interface{}) error {
	return nil
}

func (i *IPC) RPCSubscribe(args *SubArgs, reply *SubReply) error {
	i.bus.Subscribe(args.Topic, func(data any) error {
		remote, err := url.Parse(args.RemoteURL)
		if err != nil {
			log.Fatal("parse url error: ", err)
			return err
		}
		client, err := rpc.DialHTTP(remote.Scheme, remote.Host)
		if err != nil {
			log.Fatal("Client connection error: ", err)
			return err
		}

		reply := &PubReply{}
		err = client.Call("IPC.RPCPublish", &PubArgs{
			Topic: args.Topic,
			data:  data,
		}, &reply)
		if err != nil {
			log.Fatal("Client invocation error: ", err)
			return err
		}
		return nil
	})
	return nil
}

func (i *IPC) RPCPublish(args *PubArgs, reply *PubReply) error {
	i.bus.Publish(args.Topic, args.data)
	return nil
}
