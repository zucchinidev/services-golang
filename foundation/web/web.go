package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
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
}

// NewApp creates a new App value that contains the information for the HTTP server.
func NewApp(shutdown chan os.Signal) *App {
	return &App{
		ServeMux: http.NewServeMux(),
		shutdown: shutdown,
	}
}

func (a *App) HandleFunc(pattern string, handler Handler) {

	h := func(w http.ResponseWriter, r *http.Request) {

		// Put any code here that needs to happen before the handler is called.

		if err := handler(r.Context(), w, r); err != nil {
			// Temporary logging.
			fmt.Println("web.HandleFunc: %w", err)
			return
		}

		// Put any code here that needs to happen after the handler is called.
	}

	a.ServeMux.HandleFunc(pattern, h)
}
