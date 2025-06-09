// Package mid provides protocol-agnostic middleware for the application.
// This package focuses on business logic middleware that is independent of
// delivery mechanisms (HTTP, gRPC, etc.). It provides the core middleware
// functionality that can be adapted by different protocol-specific layers.
//
// This is one of the few packages where the package name describes what it contains
// (contains middleware) rather than what it does, which is an acceptable exception.

package mid

import "context"

// Handler is a function that takes a context and returns an error.
// It is used to execute any middleware function
// independent of the protocol.
// Do not confuse this with the web.Handler type.
type Handler func(context.Context) error
