package service

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
)

type OperationType string

const (
	operationTypeFill  OperationType = "fill"
	operationTypeEmpty OperationType = "empty"
	operationTypePour  OperationType = "pour"

	xJugTag = "x"
	yJugTag = "y"
)

type Operation struct {
	OperationType  `json:"operation,omitempty"`
	Jug            *string `json:"jug,omitempty"`
	JugOrigin      *string `json:"jug_origin,omitempty"`
	JugDestination *string `json:"jug_destination,omitempty"`
	WaterAmount    int     `json:"amount,omitempty"`
	Step           int     `json:"step,omitempty"`
	Description    string  `json:"description,omitempty"`
}

type RiddleResponse struct {
	Operations []Operation `json:"operations,omitempty"`
	Jug        string      `json:"jug,omitempty"`
	TotalSteps int         `json:"total_steps,omitempty"`
}

func (s *service) Riddle(x, y, z int) (*RiddleResponse, *AppError) {
	operations, jug, err := s.getOperations(x, y, z)
	if err != nil {
		return nil, err
	}

	return &RiddleResponse{
		Operations: operations,
		Jug:        jug,
		TotalSteps: len(operations),
	}, nil
}

func (s *service) getOperations(x, y, z int) ([]Operation, string, *AppError) {
	smallerJug := x
	biggerJug := y

	smallerJugTag := xJugTag
	biggerJugTag := yJugTag

	// Assumming that x < y, otherwise rename them
	if x > y {
		t := x
		smallerJug = y
		biggerJug = t
		smallerJugTag = yJugTag
		biggerJugTag = xJugTag
	}

	if z > biggerJug {
		return nil, "", &AppError{
			Error:   fmt.Errorf("can't measure %d if it's bigger than jugs for %d and %d", z, smallerJug, biggerJug),
			Message: "invalid parameters",
			Code:    http.StatusBadRequest,
		}
	}

	// If gcd of smaller jug and bigger jug does not divide z, then solution is not possible
	calculatedGcd := gcd(smallerJug, biggerJug)
	if (z % calculatedGcd) != 0 {
		return nil, "", &AppError{
			Error: fmt.Errorf("there is no solution to measure %d with jugs with %d and %d", z, smallerJug, biggerJug),
			Message: "invalid parameters",
			Code:  http.StatusBadRequest,
		}
	}


	// Test two possible scenarios
	wg := sync.WaitGroup{}
	wg.Add(2)

	// a) Water of bigger jug is poured into smaller jug
	var firstSolutionOperations []Operation
	var firstSolutionJug string

	go func() {
		firstSolutionOperations, firstSolutionJug = s.pour(biggerJug, biggerJugTag, smallerJug, smallerJugTag, z)
		wg.Done()
	}()


	// b) Water of smaller jug is poured into bigger jug
	var secondSolutionOperations []Operation
	var secondSolutionJug string

	go func() {
		secondSolutionOperations, secondSolutionJug = s.pour(smallerJug, smallerJugTag, biggerJug, biggerJugTag, z)
		wg.Done()
	}()

	wg.Wait()

	if len(firstSolutionOperations) < len(secondSolutionOperations) {
		return firstSolutionOperations, firstSolutionJug, nil
	}
	return secondSolutionOperations, secondSolutionJug, nil
}

/*
 pour returns all the operations required to measure z amount of water by constantly pouring water from jug with name
      jug1Tag into jug with name jug2Tag

 []Operation: contains the list of operations
 string:      contains the tag of the jug with the amount of water requested
*/

func (s *service) pour(jug1Cap int, jug1Tag string, jug2Cap int, jug2Tag string, z int) (ops []Operation,
	jugTag string) {
	var operations []Operation

	jug1 := jug1Cap
	jug2 := 0

	step := 1
	operations = append(operations, Operation{
		OperationType: operationTypeFill,
		Jug:           aws.String(jug1Tag),
		WaterAmount:   jug1Cap,
		Description:   fmt.Sprintf("filling jug %s with %d capacity", jug1Tag, jug1Cap),
		Step:          step,
	})

	// Break the loop when either of the two jugs has z water
	for jug1 != z && jug2 != z {
		// Find the maximum amount that can be poured
		temp := min(jug1, jug2Cap-jug2)

		// Pour "temp" liters from "jug1" to "jug2"
		jug2 += temp
		jug1 -= temp

		step++
		operations = append(operations, Operation{
			OperationType:  operationTypePour,
			JugOrigin:      aws.String(jug1Tag),
			JugDestination: aws.String(jug2Tag),
			WaterAmount:    temp,
			Description:    fmt.Sprintf("pouring water from jug %s to %s", jug1Tag, jug2Tag),
			Step:           step,
		})

		if jug1 == z {
			jugTag = jug1Tag
			break
		}

		if jug2 == z {
			jugTag = jug2Tag
			break
		}

		// If first jug becomes empty, fill it
		if jug1 == 0 {
			jug1 = jug1Cap
			step++
			operations = append(operations, Operation{
				OperationType: operationTypeFill,
				Jug:           aws.String(jug1Tag),
				WaterAmount:   jug1Cap,
				Description:   fmt.Sprintf("filling jug %s with %d capacity", jug1Tag, jug1Cap),
				Step:          step,
			})
		}

		// If second jug becomes full, empty it
		if jug2 == jug2Cap {
			jug2 = 0
			step++
			operations = append(operations, Operation{
				OperationType: operationTypeEmpty,
				Jug:           aws.String(jug2Tag),
				WaterAmount:   jug2Cap,
				Description:   fmt.Sprintf("emptying jug %s with %d capacity", jug2Tag, jug2Cap),
				Step:          step,
			})
		}
	}

	ops = operations
	return
}
