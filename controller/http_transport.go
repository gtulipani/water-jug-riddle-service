package controller

import (
	"fmt"
	"net/http"
	"water-jug-riddle-service/service"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	rice "github.com/GeertJohan/go.rice"
)

const (
	apiResource    = "api"
	v1Resource     = "v1"
	healthResource = "health"
	riddleResource = "riddle"
	staticResource = "static"
	wildcardResource     = "*"
)

var (
	healthEndpoint = fmt.Sprintf("/%s/%s/%s", apiResource, v1Resource, healthResource)
	riddleEndpoint = fmt.Sprintf("/%s/%s/%s", apiResource, v1Resource, riddleResource)
	staticEndpint = fmt.Sprintf("/%s/%s", staticResource, wildcardResource)
	rootEndpint   = "/"
)

// NewHandler: create handlers
func NewHandler(box *rice.Box, svc service.Service) http.Handler {
	r := chi.NewRouter()

	r.Handle(staticEndpint, http.FileServer(box.HTTPBox()))
	r.HandleFunc(rootEndpint, serve(box, svc))

	r.Group(func(r chi.Router) {
		r.Use(middleware.Recoverer, middleware.StripSlashes, middleware.Logger)

		r.Get(healthEndpoint, health(svc))
		r.Get(riddleEndpoint, riddle(svc))
	})

	return r
}
