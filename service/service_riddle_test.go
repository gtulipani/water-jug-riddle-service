package service

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
)

func TestService_Riddle(t *testing.T) {
	type args struct {
		x int
		y int
		z int
	}
	type want struct {
		output    *RiddleResponse
		outputErr *AppError
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "z is bigger than x and y",
			args: args{
				x: 1,
				y: 2,
				z: 3,
			},
			want: want{
				outputErr: &AppError{
					Error:   fmt.Errorf("can't measure %d if it's bigger than jugs for %d and %d",3, 1, 2),
					Message: "invalid parameters",
					Code:    http.StatusBadRequest,
				},
			},
		},
		{
			name: "gcd of smaller jug and bigger jug doesn't divide z, i.e. x and y are multiples and z is not",
			args: args{
				x: 2,
				y: 4,
				z: 3,
			},
			want: want{
				outputErr: &AppError{
					Error: fmt.Errorf("there is no solution to measure %d with jugs with %d and %d", 3, 2, 4),
					Message: "invalid parameters",
					Code:  http.StatusBadRequest,
				},
			},
		},
		{
			name: "success with x = 3, y = 5 and z = 4",
			args: args{
				x: 3,
				y: 5,
				z: 4,
			},
			want: want{
				output: &RiddleResponse{
					Operations: []Operation{
						{
							OperationType: operationTypeFill,
							Jug: aws.String(yJugTag),
							WaterAmount: 5,
							Step: 1,
							Description: fmt.Sprintf("filling jug %s with 5 capacity", yJugTag),
						},
						{
							OperationType: operationTypePour,
							JugOrigin: aws.String(yJugTag),
							JugDestination: aws.String(xJugTag),
							WaterAmount: 3,
							Step: 2,
							Description: fmt.Sprintf("pouring water from jug %s to %s", yJugTag, xJugTag),
						},
						{
							OperationType: operationTypeEmpty,
							Jug: aws.String(xJugTag),
							WaterAmount: 3,
							Step: 3,
							Description: fmt.Sprintf("emptying jug %s with 3 capacity", xJugTag),
						},
						{
							OperationType: operationTypePour,
							JugOrigin: aws.String(yJugTag),
							JugDestination: aws.String(xJugTag),
							WaterAmount: 2,
							Step: 4,
							Description: fmt.Sprintf("pouring water from jug %s to %s", yJugTag, xJugTag),
						},
						{
							OperationType: operationTypeFill,
							Jug: aws.String(yJugTag),
							WaterAmount: 5,
							Step: 5,
							Description: fmt.Sprintf("filling jug %s with 5 capacity", yJugTag),
						},
						{
							OperationType: operationTypePour,
							JugOrigin: aws.String(yJugTag),
							JugDestination: aws.String(xJugTag),
							WaterAmount: 1,
							Step: 6,
							Description: fmt.Sprintf("pouring water from jug %s to %s", yJugTag, xJugTag),
						},
					},
					Jug: yJugTag,
					TotalSteps: 6,
				},
			},
		},
		{
			name: "success with y = 3, x = 5 and z = 4",
			args: args{
				y: 3,
				x: 5,
				z: 4,
			},
			want: want{
				output: &RiddleResponse{
					Operations: []Operation{
						{
							OperationType: operationTypeFill,
							Jug: aws.String(xJugTag),
							WaterAmount: 5,
							Step: 1,
							Description: fmt.Sprintf("filling jug %s with 5 capacity", xJugTag),
						},
						{
							OperationType: operationTypePour,
							JugOrigin: aws.String(xJugTag),
							JugDestination: aws.String(yJugTag),
							WaterAmount: 3,
							Step: 2,
							Description: fmt.Sprintf("pouring water from jug %s to %s", xJugTag, yJugTag),
						},
						{
							OperationType: operationTypeEmpty,
							Jug: aws.String(yJugTag),
							WaterAmount: 3,
							Step: 3,
							Description: fmt.Sprintf("emptying jug %s with 3 capacity", yJugTag),
						},
						{
							OperationType: operationTypePour,
							JugOrigin: aws.String(xJugTag),
							JugDestination: aws.String(yJugTag),
							WaterAmount: 2,
							Step: 4,
							Description: fmt.Sprintf("pouring water from jug %s to %s", xJugTag, yJugTag),
						},
						{
							OperationType: operationTypeFill,
							Jug: aws.String(xJugTag),
							WaterAmount: 5,
							Step: 5,
							Description: fmt.Sprintf("filling jug %s with 5 capacity", xJugTag),
						},
						{
							OperationType: operationTypePour,
							JugOrigin: aws.String(xJugTag),
							JugDestination: aws.String(yJugTag),
							WaterAmount: 1,
							Step: 6,
							Description: fmt.Sprintf("pouring water from jug %s to %s", xJugTag, yJugTag),
						},
					},
					Jug: xJugTag,
					TotalSteps: 6,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &service{}
			output, outputErr := svc.Riddle(tt.args.x, tt.args.y, tt.args.z)

			a := assert.New(t)

			if tt.want.output == nil {
				a.Nil(output)
			} else {
				a.Equal(tt.want.output.Operations, output.Operations)
				a.Equal(tt.want.output.TotalSteps, output.TotalSteps)
				a.Equal(tt.want.output.Jug, output.Jug)
			}
			a.Equal(tt.want.outputErr, outputErr)
		})
	}
}
