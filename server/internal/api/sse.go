package api

import (
	"encoding/json"

	"fmt"
	"net/http"
)

// WriteSSE marshals data as a single SSE event and flushes immediately.
func WriteSSE(w http.ResponseWriter, flusher http.Flusher, event string, data any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, payload)

	flusher.Flush()
	return nil
}
