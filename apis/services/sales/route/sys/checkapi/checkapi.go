// Package checkapi provides support to check the API health.
package checkapi

import (
	"context"
	"encoding/json"
	"net/http"
)

func liveness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	return json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	return json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
