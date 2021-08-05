package controller

import (
	"net/http"
	"time"
	"water-jug-riddle-service/service"

	rice "github.com/GeertJohan/go.rice"
)

func serve(app *rice.Box, svc service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		indexFile, err := app.Open("index.html")
		if err != nil {
			encodeHTTPError(&service.AppError{
				Error: err,
				Message: "error opening index.html file",
				Code: http.StatusInternalServerError,
			}, w)
			return
		}

		http.ServeContent(w, r, "index.html", time.Time{}, indexFile)
	}
}