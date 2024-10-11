package api

import (
	"net/http"

	"github.com/go-chi/render"
)

type response struct {
	status int         `json:"-"`
	Data   interface{} `json:"data,omitempty"`
	Meta   interface{} `json:"meta,omitempty"`
	Error  string      `json:"error,omitempty"`
}

func Response(status int) *response {
	return &response{
		status: status,
	}
}

func (r *response) WithBody(data, meta interface{}) *response {
	if r != nil {
		r.Data = data
		r.Meta = meta
	}

	return r
}

func (r *response) WithError(err error) *response {
	if r != nil {
		r.Error = err.Error()
	}

	return r
}

func (r *response) Send(writer http.ResponseWriter, request *http.Request) {
	if r != nil {
		render.Status(request, r.status)
		render.JSON(writer, request, r)
	}
}
