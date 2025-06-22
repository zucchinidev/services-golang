package mid

import (
	"context"

	"github.com/zucchini/services-golang/app/api/metrics"
)

func Metrics(ctx context.Context, handler Handler) error {

	ctx = metrics.Set(ctx)
	err := handler(ctx)

	requests := metrics.AddRequests(ctx)

	if requests%5000 == 0 {
		// This is a hack to get the number of goroutines.
		// We should use a more robust way to get the number of goroutines.
		// My recommendation is to call this function into the liveness handler.
		// Or maybe use the uber library: https://github.com/uber-go/tally
		metrics.AddGoroutines(ctx)
	}

	if err != nil {
		metrics.AddErrors(ctx)
	}

	return err
}
