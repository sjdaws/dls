package web

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

// CopyBody creates a copy of the request body without modifying the original
func CopyBody(request *http.Request) io.ReadCloser {
	buffer, _ := io.ReadAll(request.Body)
	request.Body = io.NopCloser(bytes.NewBuffer(buffer))

	return io.NopCloser(bytes.NewBuffer(buffer))
}

// DecodeBody decodes the request body into a struct and fixes common issues
func DecodeBody(request *http.Request, into any) error {
	err := json.NewDecoder(CopyBody(request)).Decode(&into)
	if err == nil {
		return nil
	}

	// Attempt to fix
	buffer, _ := io.ReadAll(CopyBody(request))
	strBuffer := string(buffer)

	// Invalid mac address list
	if !strings.Contains(strBuffer, `"mac_address_list":["`) {
		strBuffer = strings.Replace(strBuffer, `"mac_address_list":[`, `"mac_address_list":["`, 1)
		request.Body = io.NopCloser(bytes.NewBuffer([]byte(strBuffer)))
		err = json.NewDecoder(CopyBody(request)).Decode(&into)
		if err != nil {
			return err
		}
	}

	return nil
}

// getRemoteAddress will get the originating IP(s)
func getRemoteAddress(request *http.Request) string {
	forwardedFor := request.Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		return forwardedFor
	}

	return request.RemoteAddr
}
