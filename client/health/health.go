package health

import (
    "net/http"

    "github.com/sjdaws/dls/internal/web"
)

// GetHealth returns "OK" if the service is up
func GetHealth(response http.ResponseWriter, request *http.Request) {
    web.Respond(request, response, "text/plain", "OK")
}
