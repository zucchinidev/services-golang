package authclient

import (
	"github.com/google/uuid"
	"github.com/zucchini/services-golang/business/api/auth"
)

type Error struct {
	Message string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}

// Authorize defines the information required to perform an authorization
type Authorize struct {
	Claims auth.Claims
	UserID uuid.UUID
	Rule   string
}

// AuthenticateResp defines the information that will be received on authenticate
type AuthenticateResp struct {
	UserID uuid.UUID
	Claims auth.Claims
}
