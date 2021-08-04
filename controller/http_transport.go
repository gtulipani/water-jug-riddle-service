package controller

import (
	"fmt"
	"net/http"
	"water-jug-riddle-service/service"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	apiResource    = "api"
	v1Resource     = "v1"
	healthResource = "health"
	riddleResource = "riddle"
)

var (
	healthEndpoint = fmt.Sprintf("/%s/%s/%s", apiResource, v1Resource, healthResource)
	riddleEndpoint = fmt.Sprintf("/%s/%s/%s", apiResource, v1Resource, riddleResource)
)

// NewHandler: create handlers
func NewHandler(svc service.Service) http.Handler {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(middleware.Recoverer, middleware.StripSlashes, middleware.Logger)

		r.Get(healthEndpoint, health(svc))
		r.Get(riddleEndpoint, riddle(svc))
	})

	return r
}
