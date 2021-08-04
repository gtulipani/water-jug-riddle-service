package controller

import (
	"net/http"
	"water-jug-riddle-service/service"
)

func health(svc service.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		response := svc.Health()

		if err := encodeHTTPResponse(w, response); err != nil {
			encodeHTTPError(err, w)
		}
	}
}