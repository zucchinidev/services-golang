package web

import (
	"context"
	"time"
)

// ctxKey is a type for the context key.
// It is used to store the Values struct in the context.
// We do not want to use this directly in the code, so we will create an API for it.
type ctxKey int

const key ctxKey = 1

// Values represent state for each request.
type Values struct {
	TraceID    string
	Now        time.Time
	StatusCode int
}

// GetValues returns the Values struct from the context.
// If the context does not contain a Values struct, it creates a new one.
func GetValues(ctx context.Context) *Values {

	v, ok := ctx.Value(key).(*Values)

	// If the context does not contain a Values struct, create a new one.
	// We will avoid errors when using Values.
	if !ok {
		return &Values{
			TraceID:    "00000000-0000-0000-0000-000000000000",
			Now:        time.Now(),
			StatusCode: 200,
		}
	}

	return v
}

func GetTraceID(ctx context.Context) string {
	return GetValues(ctx).TraceID
}

func GetTime(ctx context.Context) time.Time {
	return GetValues(ctx).Now
}

func GetStatusCode(ctx context.Context) int {
	return GetValues(ctx).StatusCode
}

func setStatusCode(ctx context.Context, statusCode int) {
	v, ok := ctx.Value(key).(*Values)
	if !ok {
		return
	}

	v.StatusCode = statusCode
}

func setValues(ctx context.Context, v *Values) context.Context {
	return context.WithValue(ctx, key, v)
}
