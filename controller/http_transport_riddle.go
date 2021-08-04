package controller

import (
	"errors"
	"net/http"
	"strconv"
	"water-jug-riddle-service/service"
)

const (
	xQueryParam = "x"
	yQueryParam = "y"
	zQueryParam = "z"
)

type RiddleRequest struct {
	X int `json:"x,omitempty"`
	Y int `json:"y,omitempty"`
	Z int `json:"z,omitempty"`
}

func riddle(svc service.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := decodeRiddleRequest(r)
		if err != nil {
			encodeHTTPError(err, w)
			return
		}

		response, err := svc.Riddle(req.X, req.Y, req.Z)
		if err != nil {
			encodeHTTPError(err, w)
			return
		}

		if err := encodeHTTPResponse(w, response); err != nil {
			encodeHTTPError(err, w)
		}
	}
}

func decodeRiddleRequest(r *http.Request) (*RiddleRequest, *service.AppError) {
	x, err := getIntegerQueryParam(r, xQueryParam)
	if err != nil {
		return nil, &service.AppError{
			Error: err,
			Message: "invalid parameters",
			Code: http.StatusBadRequest,
		}
	}

	y, err := getIntegerQueryParam(r, yQueryParam)
	if err != nil {
		return nil, &service.AppError{
			Error: err,
			Message: "invalid parameters",
			Code: http.StatusBadRequest,
		}
	}

	z, err := getIntegerQueryParam(r, zQueryParam)
	if err != nil {
		return nil, &service.AppError{
			Error: err,
			Message: "invalid parameters",
			Code: http.StatusBadRequest,
		}
	}

	if valid := validateRiddleRequest(x, y, z); !valid {
		return nil, &service.AppError{
			Error: errors.New("every param must be a positive integer"),
			Message: "invalid parameters",
			Code: http.StatusBadRequest,
		}
	}

	return &RiddleRequest{
		X: x,
		Y: y,
		Z: z,
	}, nil
}

func getIntegerQueryParam(r *http.Request, param string) (int, error) {
	stringValue := r.URL.Query().Get(param)
	if stringValue == "" {
		return 0, errors.New("every param must be a positive integer")
	}
	intValue, err :=  strconv.Atoi(stringValue)
	if err != nil {
		return 0, errors.New("value is not integer")
	}
	return intValue, nil
}

func validateRiddleRequest(x, y, z int) bool {
	return x > 0 && y > 0 && z > 0
}
