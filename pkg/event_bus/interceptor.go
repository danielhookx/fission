package eventbus

import (
	"context"
	"fmt"
)

type Handler func(ctx context.Context, topic string, data any) error
type Interceptor func(ctx context.Context, topic string, data any, handler Handler) error

func eventbusInterceptors(e *EventBus) {
	// Prepend opts.unaryInt to the chaining interceptors if it exists, since unaryInt will
	// be executed before any other chained interceptors.
	interceptors := e.opts.interceptors
	var chainedInt Interceptor
	if len(interceptors) == 0 {
		chainedInt = nil
	} else if len(interceptors) == 1 {
		chainedInt = interceptors[0]
	} else {
		chainedInt = chainInterceptors(interceptors)
	}
	// e.interceptor = chainedInt
	fmt.Println(chainedInt)
}

func chainInterceptors(interceptors []Interceptor) Interceptor {
	return func(ctx context.Context, topic string, data any, handler Handler) error {
		// the struct ensures the variables are allocated together, rather than separately, since we
		// know they should be garbage collected together. This saves 1 allocation and decreases
		// time/call by about 10% on the microbenchmark.
		var state struct {
			i    int
			next Handler
		}
		state.next = func(ctx context.Context, topic string, data any) error {
			if state.i == len(interceptors)-1 {
				return interceptors[state.i](ctx, topic, data, handler)
			}
			state.i++
			return interceptors[state.i-1](ctx, topic, data, state.next)
		}
		return state.next(ctx, topic, data)
	}
}
