// Package mid provides protocol-agnostic middleware for the application.
// This package focuses on business logic middleware that is independent of
// delivery mechanisms (HTTP, gRPC, etc.). It provides the core middleware
// functionality that can be adapted by different protocol-specific layers.
//
// This is one of the few packages where the package name describes what it contains
// (contains middleware) rather than what it does, which is an acceptable exception.

package mid

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/zucchini/services-golang/business/api/auth"
)

// Handler is a function that takes a context and returns an error.
// It is used to execute any middleware function
// independent of the protocol.
// Do not confuse this with the web.Handler type.
type Handler func(context.Context) error

type ctxKey int

const (
	claimKey ctxKey = iota + 1
	userIDKey
)

func setClaims(ctx context.Context, claims auth.Claims) context.Context {
	return context.WithValue(ctx, claimKey, claims)
}

// GetClaims returns the claims from the context.
func GetClaims(ctx context.Context) auth.Claims {
	v, ok := ctx.Value(claimKey).(auth.Claims)
	if !ok {
		return auth.Claims{}
	}
	return v
}

func setUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// GetUserID returns the claims from the context.
func GetUserID(ctx context.Context) (uuid.UUID, error) {
	v, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, errors.New("user id not found in context")
	}

	return v, nil
}
