package web

import (
	"context"
	"encoding/json"
	"net/http"
)

// Respond is a helper function that writes data to the response writer.
// We do not need to use json around the project, we can easily change to protobuf or any other format
// by changing the implementation of this function.
func Respond(ctx context.Context, w http.ResponseWriter, data any, statusCode int) error {
	setStatusCode(ctx, statusCode)

	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if _, err := w.Write(b); err != nil {
		return err
	}

	return nil
}
