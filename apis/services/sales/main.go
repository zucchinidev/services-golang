package main

import (
	"context"
	"os"
	"runtime"

	"github.com/zucchini/services-golang/foundation/logger"
)

func main() {
	var log *logger.Logger

	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "******* SEND ALERT *******")
		},
	}

	traceIDFn := func(ctx context.Context) string {
		return ""
	}

	log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "SALES", traceIDFn, events)

	// -------------------------------------------------------------------------

	ctx := context.Background()

	if err := run(ctx, log); err != nil {
		log.Error(ctx, "startup", "msg", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {

	// -------------------------------------------------------------------------
	// GOMAXPROCS

	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// -------------------------------------------------------------------------
	// Configuration

	return nil
}