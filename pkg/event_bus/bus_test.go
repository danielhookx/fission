package eventbus

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func a(name string) {
	fmt.Printf("sub1 -- %s\n", name)
}

func b(name string) {
	fmt.Printf("sub2 -- %s\n", name)
}

func c(name string) {
	fmt.Printf("sub3 -- %s\n", name)
}

func BenchmarkSubPub(b *testing.B) {
	e := NewEventBus()
	topic := "testpub1"
	b.RunParallel(func(pb *testing.PB) {
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			for pb.Next() {
				e.Subscribe(topic, func(name string) {})
			}
		}()

		go func() {
			defer wg.Done()
			for pb.Next() {
				e.Publish(topic, "jack")
			}
		}()

		wg.Wait()
	})
}

func TestNamedSubscribe(t *testing.T) {
	e := NewEventBus()
	topic := "testpub1"
	e.Subscribe(topic, a)
	e.Subscribe(topic, b)
	e.Subscribe(topic, c)
	e.Publish(topic, "jack")
	e.Unsubscribe(topic, b)
	e.Publish(topic, "jack2")
}

func TestAnonymousSubscribe(t *testing.T) {
	e := NewEventBus()
	topic := "testpub1"
	e.Subscribe(topic, func(name string) {
		fmt.Printf("sub1 -- %s\n", name)
	})
	e.Subscribe(topic, func(name string) {
		fmt.Printf("sub2 -- %s\n", name)
	})
	e.Subscribe(topic, func(name string) {
		fmt.Printf("sub3 -- %s\n", name)
	})
	e.Publish(topic, "jack")
}

func TestBlockSubscribe(t *testing.T) {
	e := NewEventBus()
	topic := "testpub1"
	e.Subscribe(topic, func(name string) {
		fmt.Printf("sub1 -- %s\n", name)
		time.Sleep(time.Second * 1)
	})
	e.Subscribe(topic, func(name string) {
		fmt.Printf("sub2 -- %s\n", name)
		time.Sleep(time.Second * 2)
	})
	e.Subscribe(topic, func(name string) {
		fmt.Printf("sub3 -- %s\n", name)
		time.Sleep(time.Second * 3)
	})
	e.Publish(topic, "jack")
}

type testA struct {
	name string
}

func TestSubscribeParamIsolation(t *testing.T) {
	e := NewEventBus()
	topic := "testpub1"
	var s1 = testA{
		name: "jack",
	}
	e.Subscribe(topic, func(v *testA) {
		fmt.Printf("sub1 -- %s\n", v.name)
		v.name = "Lee"
	})
	e.Subscribe(topic, func(v *testA) {
		fmt.Printf("sub2 -- %s\n", v.name)
		v.name = "Danny"
	})
	e.Subscribe(topic, func(v *testA) {
		fmt.Printf("sub3 -- %s\n", v.name)
		v.name = "Jay"
	})
	e.Publish(topic, &s1)
}
