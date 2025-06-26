package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aliyun/fc-runtime-go-sdk/events"
	"github.com/danvei233/fchandler2go"
	"github.com/gin-gonic/gin"
	"net/http"
)

//this is a test file

func main() {
	router := gin.Default()
	router.Any("/*all", handler2gin.T(HandleRequest))

	router.Run(":8080")
}

//	type StructEvent struct {
//		Key string `json:"key"`
//	}
//
//	func HandleRequest(ctx context.Context, event StructEvent) (string, error) {
//		return fmt.Sprintf("hello, %s!", event.Key), nil
//	}
type HTTPTriggerEvent events.HTTPTriggerEvent
type HTTPTriggerResponse events.HTTPTriggerResponse

func (h HTTPTriggerEvent) String() string {
	jsonBytes, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

func NewHTTPTriggerResponse(statusCode int) *HTTPTriggerResponse {
	return &HTTPTriggerResponse{StatusCode: statusCode}
}

func (h *HTTPTriggerResponse) String() string {
	jsonBytes, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

func (h *HTTPTriggerResponse) WithStatusCode(statusCode int) *HTTPTriggerResponse {
	h.StatusCode = statusCode
	return h
}

func (h *HTTPTriggerResponse) WithHeaders(headers map[string]string) *HTTPTriggerResponse {
	h.Headers = headers
	return h
}

func (h *HTTPTriggerResponse) WithIsBase64Encoded(isBase64Encoded bool) *HTTPTriggerResponse {
	h.IsBase64Encoded = isBase64Encoded
	return h
}

func (h *HTTPTriggerResponse) WithBody(body string) *HTTPTriggerResponse {
	h.Body = body
	return h
}

func GetRawRequestEvent(event []byte) (*HTTPTriggerResponse, error) {
	fmt.Printf("raw event: %s\n", string(event))
	return NewHTTPTriggerResponse(http.StatusOK).WithBody(string(event)), nil
}

func HandleRequest(event HTTPTriggerEvent) (*HTTPTriggerResponse, error) {
	fmt.Printf("event: %v\n", event)
	if event.Body == nil {
		return NewHTTPTriggerResponse(http.StatusBadRequest).
			WithBody(fmt.Sprintf("the request did not come from an HTTP Trigger, event: %v", event)), nil
	}

	reqBody := *event.Body
	if event.IsBase64Encoded != nil && *event.IsBase64Encoded {
		decodedByte, err := base64.StdEncoding.DecodeString(*event.Body)
		if err != nil {
			return NewHTTPTriggerResponse(http.StatusBadRequest).
				WithBody(fmt.Sprintf("HTTP Trigger body is not base64 encoded, err: %v", err)), nil
		}
		reqBody = string(decodedByte)
	} else {
		return NewHTTPTriggerResponse(http.StatusBadRequest).
			WithBody("HELLO WORLD"), nil

	}

	return NewHTTPTriggerResponse(http.StatusOK).WithBody(reqBody), nil
}
