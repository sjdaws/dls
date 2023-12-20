package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/sjdaws/dls/internal/global"
	"github.com/sjdaws/dls/internal/notify"
)

type HttpError struct {
	Detail     string `json:"detail"`
	StatusCode int    `json:"status"`
}

// Error returns an error response
func Error(request *http.Request, response http.ResponseWriter, message string, err error, httpError *HttpError) {
	log.Printf("ERROR: %s", fmt.Sprintf("%s via %s: %v", message, request.URL.Path, err))
	go notify.Message(fmt.Sprintf("Error encountered.\n\nMessage: %s\nError: %v\n\nIP: %s\nMethod: %s\nURL: %s", message, err, getRemoteAddress(request), request.Method, request.URL.Path))

	if httpError == nil {
		httpError = &HttpError{
			Detail:     "Internal error",
			StatusCode: http.StatusInternalServerError,
		}
	}

	body, _ := json.Marshal(httpError)

	send(request, response, "application/json", string(body), httpError.StatusCode)
}

func Respond(request *http.Request, response http.ResponseWriter, contentType string, message string, parameters ...any) {
	send(request, response, contentType, fmt.Sprintf(message, parameters...), http.StatusOK)
}

func send(request *http.Request, response http.ResponseWriter, contentType string, body string, statusCode int) {
	if global.Debug {
		buffer, _ := io.ReadAll(request.Body)
		log.Printf("\nIP: %s\nMethod: %s\nURL: %s\nAuthorization: %s\nBody: %s\nResponse: %s\n", getRemoteAddress(request), request.Method, request.URL.Path, request.Header.Get("authorization"), buffer, body)
		request.Body = io.NopCloser(bytes.NewBuffer(buffer))
	}

	response.Header().Add("Content-Length", fmt.Sprintf("%d", len(body)))
	response.WriteHeader(statusCode)
	_, _ = io.WriteString(response, body)
}
