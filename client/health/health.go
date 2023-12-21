package health

import (
	"fmt"
	"io"
	"net/http"
)

// GetHealth returns "OK" if the service is up
func GetHealth(response http.ResponseWriter, request *http.Request) {
	body := "OK"

	// Write response directly to bypass debug logging
	response.Header().Add("Content-Length", fmt.Sprintf("%d", len(body)))
	response.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(response, body)
}
