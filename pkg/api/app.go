package api

import (
	"time"

	"github.com/go-chi/chi/v5"
)

const (
	Timeout = 5 * time.Second
)

type App interface {
	RegisterPublicEndpoints(chi.Router)
	RegisterProtectedEndpoints(chi.Router)
}
