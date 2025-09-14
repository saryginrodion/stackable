package stackable

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type Response interface {
	Headers() HeadersContainer
	Body() io.Reader
	Status() int
}

type HttpResponse struct {
	body    io.Reader
	status  int
	headers HeadersContainer
}

func NewHttpResponseRaw(headers HeadersContainer, status int, body io.Reader) *HttpResponse {
	resp := HttpResponse{
		headers: headers,
		status:  status,
		body:    body,
	}

	return &resp
}

func NewHttpResponse(status int, contentType string, body string) *HttpResponse {
	headers := NewHeadersContainer()
	headers.Set("Content-Type", contentType)
	resp := HttpResponse{
		headers: headers,
		status:  status,
		body:    strings.NewReader(body),
	}

	return &resp
}

func (r *HttpResponse) SetHeaders(newHeaders HeadersContainer) {
	r.headers = newHeaders
}

func (r *HttpResponse) Headers() HeadersContainer {
	return r.headers
}

func (r *HttpResponse) Body() io.Reader {
	return r.body
}

func (r *HttpResponse) Status() int {
	return r.status
}

func JsonResponse(status int, data any) (*HttpResponse, error) {
	jsonBytes, err := json.Marshal(data)

	if err != nil {
		logrus.Error("Failed to serialize json. ", err)
		return nil, HttpError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to serialise JSON. Error: " + err.Error(),
		}
	}

	resp := NewHttpResponse(
		status,
		"application/json",
		string(jsonBytes),
	)

	return resp, nil
}
