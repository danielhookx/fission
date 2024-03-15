package eventbus

import (
	"context"
	"fmt"
	"reflect"

	"github.com/danielhookx/fission"
)

func Func(fn reflect.Value) context.Context {
	return context.WithValue(context.Background(), "eventbus", fn)
}

func FromCtx(ctx context.Context) reflect.Value {
	fn := ctx.Value("eventbus")
	return fn.(reflect.Value)
}

type BusSubscriber interface {
	Subscribe(topic string, fn interface{}) error
	SubscribeSync(topic string, fn interface{}) error
	Unsubscribe(topic string, handler interface{}) error
}

type EventBus struct {
	rm *fission.RouteManager
	pm *fission.PlatformManager
}

func NewEventBus() *EventBus {
	e := &EventBus{
		rm: fission.NewRouteManager(),
	}
	e.pm = fission.NewPlatformManager(e.CreateEventBusAsyncDist)
	return e
}

func (bus *EventBus) CreateEventBusSyncDist(ctx context.Context, key any) fission.Distribution {
	fn := FromCtx(ctx)
	return NewSyncDistribution(fn)
}

func (bus *EventBus) CreateEventBusAsyncDist(ctx context.Context, key any) fission.Distribution {
	fn := FromCtx(ctx)
	return NewAsyncDistribution(fn)
}

// Wrapper function that transforms a function into a comparable interface.
func functionWrapper(f interface{}) interface{} {
	return reflect.ValueOf(f).Pointer()
}

func (bus *EventBus) Subscribe(topic string, fn interface{}) error {
	fnType := reflect.TypeOf(fn)
	if !(fnType.Kind() == reflect.Func) {
		return fmt.Errorf("%s is not of type reflect.Func", fnType.Kind())
	}

	r := bus.rm.PutRoute(topic)
	p := bus.pm.PutPlatform(Func(reflect.ValueOf(fn)), functionWrapper(fn), nil)
	r.AddPlatform(p)
	return nil
}

func (bus *EventBus) SubscribeSync(topic string, fn interface{}) error {
	fnType := reflect.TypeOf(fn)
	if !(fnType.Kind() == reflect.Func) {
		return fmt.Errorf("%s is not of type reflect.Func", fnType.Kind())
	}

	r := bus.rm.PutRoute(topic)
	p := bus.pm.PutPlatform(Func(reflect.ValueOf(fn)), functionWrapper(fn), bus.CreateEventBusSyncDist)
	r.AddPlatform(p)
	return nil
}

func (bus *EventBus) Unsubscribe(topic string, handler interface{}) error {
	fnType := reflect.TypeOf(handler)
	if !(fnType.Kind() == reflect.Func) {
		return fmt.Errorf("%s is not of type reflect.Func", fnType.Kind())
	}
	r := bus.rm.PutRoute(topic)
	r.DelPlatform(functionWrapper(handler))
	return nil
}

func (bus *EventBus) Publish(topic string, args ...interface{}) {
	r := bus.rm.PutRoute(topic)
	r.Fission(args)
	return
}

type SyncDistribution struct {
	fn reflect.Value
}

func NewSyncDistribution(fn reflect.Value) *SyncDistribution {
	return &SyncDistribution{
		fn: fn,
	}
}

func (d *SyncDistribution) Dist(data any) error {
	passedArguments := setFuncArgs(d.fn, data.([]interface{}))
	d.fn.Call(passedArguments)
	return nil
}

func (d *SyncDistribution) Close() error {
	return nil
}

type AsyncDistribution struct {
	fn reflect.Value
}

func NewAsyncDistribution(fn reflect.Value) *AsyncDistribution {
	return &AsyncDistribution{
		fn: fn,
	}
}

func (d *AsyncDistribution) Dist(data any) error {
	go func() {
		passedArguments := setFuncArgs(d.fn, data.([]interface{}))
		d.fn.Call(passedArguments)
	}()
	return nil
}

func (d *AsyncDistribution) Close() error {
	return nil
}

func setFuncArgs(fn reflect.Value, args []interface{}) []reflect.Value {
	funcType := fn.Type()
	passedArguments := make([]reflect.Value, len(args))
	for i, v := range args {
		if v == nil {
			passedArguments[i] = reflect.New(funcType.In(i)).Elem()
		} else {
			passedArguments[i] = reflect.ValueOf(v)
		}
	}
	return passedArguments
}
