package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

// Handler is a function that can handle a http request within our small little own HTTP
// framework.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App is the entrypoint to our application that configures context for HTTP handlers.
// It wraps the standard library HTTP server to provide idiomatic Go patterns like
// context as first parameter and error returns.
type App struct {
	// We embed http.ServeMux so we can use it as our router.
	// My App is not everything that http.ServeMux is.
	*http.ServeMux
	shutdown chan os.Signal
	mw       []MidHandler
}

// NewApp creates a new App value that contains the information for the HTTP server.
func NewApp(shutdown chan os.Signal, mw ...MidHandler) *App {
	return &App{
		ServeMux: http.NewServeMux(),
		shutdown: shutdown,
		mw:       mw,
	}
}

func (a *App) HandleFunc(pattern string, handler Handler, mw ...MidHandler) {
	// local middleware first
	// This allows us to for example, add an authentication middleware only to this handler.
	handler = wrapMiddleware(mw, handler)
	// general middleware afters
	handler = wrapMiddleware(a.mw, handler)

	h := func(w http.ResponseWriter, r *http.Request) {

		// Put any code here that needs to happen before the handler is called.
		// Remember, this is a fundational layer. Highly portable.
		// Do not log here or execute code specific to a handler.

		v := &Values{
			TraceID: uuid.NewString(),
			Now:     time.Now(),
		}

		ctx := setValues(r.Context(), v)

		if err := handler(ctx, w, r); err != nil {
			// Temporary logging.
			fmt.Println("web.HandleFunc: %w", err)
			return
		}

		// Put any code here that needs to happen after the handler is called.
	}

	a.ServeMux.HandleFunc(pattern, h)
}
