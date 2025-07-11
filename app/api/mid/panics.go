package mid

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/zucchini/services-golang/app/api/metrics"
)

// Panics is a middleware that recovers from panics and returns an error so it
// can be reported in Metrics and handled in Errors
func Panics(ctx context.Context, handler Handler) (err error) {
	defer func() {
		if r := recover(); r != nil {
			trace := debug.Stack()
			err = fmt.Errorf("PANIC: [%v] [%s]", r, trace)

			metrics.AddPanics(ctx)
		}
	}()

	return handler(ctx)
}
