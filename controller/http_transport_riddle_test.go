package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"water-jug-riddle-service/service"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	type args struct {
		x string
		y string
		z string
	}
	tests := []struct {
		name string
		svc  *ServiceMock
		args
		status   int
		response interface{}
		wantErr  bool
	}{
		{
			name: "missing x param",
			svc: &ServiceMock{},
			args: args{
				y: "1",
				z: "1",
			},
			response: &APIError{
				Description: "every param must be a positive integer",
				Message:     "invalid parameters",
			},
			status: http.StatusBadRequest,
			wantErr: true,
		},
		{
			name: "missing y param",
			svc: &ServiceMock{},
			args: args{
				x: "1",
				z: "1",
			},
			response: &APIError{
				Description: "every param must be a positive integer",
				Message:     "invalid parameters",
			},
			status: http.StatusBadRequest,
			wantErr: true,
		},
		{
			name: "missing z param",
			svc: &ServiceMock{},
			args: args{
				x: "1",
				y: "1",
			},
			response: &APIError{
				Description: "every param must be a positive integer",
				Message:     "invalid parameters",
			},
			status: http.StatusBadRequest,
			wantErr: true,
		},
		{
			name: "invalid x param",
			svc: &ServiceMock{},
			args: args{
				x: "a",
				y: "1",
				z: "1",
			},
			response: &APIError{
				Description: "value is not integer",
				Message:     "invalid parameters",
			},
			status: http.StatusBadRequest,
			wantErr: true,
		},
		{
			name: "invalid y param",
			svc: &ServiceMock{},
			args: args{
				x: "1",
				y: "a",
				z: "1",
			},
			response: &APIError{
				Description: "value is not integer",
				Message:     "invalid parameters",
			},
			status: http.StatusBadRequest,
			wantErr: true,
		},
		{
			name: "invalid z param",
			svc: &ServiceMock{},
			args: args{
				x: "1",
				y: "1",
				z: "a",
			},
			response: &APIError{
				Description: "value is not integer",
				Message:     "invalid parameters",
			},
			status: http.StatusBadRequest,
			wantErr: true,
		},
		{
			name: "error with service",
			svc: &ServiceMock{
				RiddleFunc: func(x int, y int, z int) (*service.RiddleResponse, *service.AppError) {
					return nil, &service.AppError{
						Error:   errors.New("some error"),
						Message: "some message",
						Code:    http.StatusInternalServerError,
					}
				},
			},
			args: args{
				x: "1",
				y: "1",
				z: "1",
			},
			response: &APIError{
				Description: "some error",
				Message:     "some message",
			},
			status: http.StatusInternalServerError,
			wantErr: true,
		},
		{
			name: "ok",
			svc: &ServiceMock{
				RiddleFunc: func(x int, y int, z int) (*service.RiddleResponse, *service.AppError) {
					return &service.RiddleResponse{
						Operations: []service.Operation{
							{
								OperationType:  "type",
								Jug:            aws.String("x"),
								JugOrigin:      aws.String("origin"),
								JugDestination: aws.String("destination"),
								WaterAmount:    1,
								Step:           1,
								Description:    "description",
							},
						},
						Jug:        "x",
						TotalSteps: 1,
					}, nil
				},
			},
			args: args{
				x: "1",
				y: "1",
				z: "1",
			},
			status: http.StatusOK,
			response: &service.RiddleResponse{
				Operations: []service.Operation{
					{
						OperationType:  "type",
						Jug:            aws.String("x"),
						JugOrigin:      aws.String("origin"),
						JugDestination: aws.String("destination"),
						WaterAmount:    1,
						Step:           1,
						Description:    "description",
					},
				},
				Jug:        "x",
				TotalSteps: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(tt.svc)
			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodGet,
				fmt.Sprintf(
					"%s?%s=%s&%s=%s&%s=%s",
					riddleEndpoint,
					xQueryParam,
					tt.args.x,
					yQueryParam,
					tt.args.y,
					zQueryParam,
					tt.args.z),
				nil)
			h.ServeHTTP(w, r)

			rawbody, _ := ioutil.ReadAll(w.Body)

			a := assert.New(t)
			a.Equal(tt.status, w.Code)

			if tt.wantErr {
				var body APIError
				if err := json.Unmarshal(rawbody, &body); err != nil {
					t.Fatalf("error unmarshalling result")
				}

				a.Equal(tt.response, &body)
			} else {
				var body service.RiddleResponse
				if err := json.Unmarshal(rawbody, &body); err != nil {
					t.Fatalf("error unmarshalling result")
				}

				a.Equal(tt.response, &body)
			}
		})
	}
}
