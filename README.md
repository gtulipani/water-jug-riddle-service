# water-jug-riddle-service

## Intro
This is a basic Backend Service that exposes an API to solve the famous Water Jug Riddle. It has been designed with 
only one API that accepts X (Jug 1 capacity), Y (Jug 2 capacity) and Z (desired capacity) and the response includes 
all the necessary "operations" that need to be performed between Jugs in order to finally get the desired water amount.
By returning a list of operations, this could be easily displayed and animated in a fancy UI.

## Instructions
The following instructions are useful to Build, Test and Run the server.

## Building and Running
### Generating mocks
Before building the project, the `mocks` need to be created, so that the Unit Tests can run successfully. In order to create them:
```
moq -out controller/mock_service.go -pkg controller ./service Service
```
Or even simpler:
```
make gen
```

### Build
After the mocks are created, the project can be easily built with:
```
go build
```

Or even simpler:
```
make build
```
This last command generates the mocks before building the code. Coverage report can also be created with the command:
```
make cover
```

### Run
#### Environment Variables
The following environment variables need to be present or must be provided in a `.env` file in the same directory as the executable:
```
HTTP_PORT=8080
```

#### Execution
Many different ways to do it:
```
go run main.go
```
Or directly:
```
./water-jug-riddle-service
```

Or even simpler
```
make run
```

## Technologies
### Go-Chi
[Go-Chi](https://github.com/go-chi/chi) has been used as the HTTP Router. It's lightweight, idiomatic and composable, therefore no further dependencies are added.

### Moq
[Moq](https://github.com/matryer/moq) has been used in order to create mocks automatically.

### Testify
[Testify](https://github.com/stretchr/testify) is used for assertions in UTs.

## Architecture
This basic Backend Service has been divided in:
- **config**: contains all the logic to retrieve environment variables. Right now only the `HTTP_PORT` but could be 
  more.
- **controller**: contains all APIs, router, decoding and encoding.
- **service**: contains all the specific business logic, including the algorithm to solve the Water Jug Riddle.

### Assumptions
- Water Jug Riddle is solvable as long as z % gcd(smallerJug, biggerJug) is not 0.
- Two mechanisms are considered possible. Pouring water every time from smallerJug into biggerJug or viceversa. Both 
  are assumed to be valid, therefore both of them are executed in parallel (using go-routines and WaitGroups).

### Improvements
- The project could be easily dockerized with a `docker-compose` or `Dockerfile`.

## CI
The project is not configured with CI yet.

## Examples
Following examples are provided to understand how the API works. 

### Using Jugs with 3 and 5 to measure 4
```
▶ curl --location --request GET 'localhost:8080/api/v1/riddle?x=3&y=5&z=4'
{
  "operations": [
    {
      "operation": "fill",
      "jug": "y",
      "amount": 5,
      "step": 1,
      "description": "filling jug y with 5 capacity"
    },
    {
      "operation": "pour",
      "jug_origin": "y",
      "jug_destination": "x",
      "amount": 3,
      "step": 2,
      "description": "pouring water from jug y to x"
    },
    {
      "operation": "empty",
      "jug": "x",
      "amount": 3,
      "step": 3,
      "description": "emptying jug x with 3 capacity"
    },
    {
      "operation": "pour",
      "jug_originΩ": "y",
      "jug_destination": "x",
      "amount": 2,
      "step": 4,
      "description": "pouring water from jug y to x"
    },
    {
      "operation": "fill",
      "jug": "y",
      "amount": 5,
      "step": 5,
      "description": "filling jug y with 5 capacity"
    },
    {
      "operation": "pour",
      "jug_origin": "y",
      "jug_destination": "x",
      "amount": 1,
      "step": 6,
      "description": "pouring water from jug y to x"
    }
  ],
  "jug": "y",
  "total_steps": 6
}
```

### Errors
#### Missing X, Y or Z parameters
```
▶ curl --location --request GET 'localhost:8080/api/v1/riddle?x=&y=5&z=4'
{
  "description": "every param must be a positive integer",
  "message": "invalid parameters"
}
```

#### Parameter is negative
```
▶ curl --location --request GET 'localhost:8080/api/v1/riddle?x=-1&y=5&z=4'
{
  "description": "every param must be a positive integer",
  "message": "invalid parameters"
}
```

#### Parameter is not a number
```
▶ curl --location --request GET 'localhost:8080/api/v1/riddle?x=a&y=5&z=4'
{
  "description": "value is not integer",
  "message": "invalid parameters"
}
```

#### Z is bigger than X and Y jugs
```
▶ curl --location --request GET 'localhost:8080/api/v1/riddle?x=3&y=5&z=6'
{
  "description": "can't measure 6 if it's bigger than jugs for 3 and 5",
  "message": "invalid parameters"
}
```

#### There is no solution, i.e. because X and Y are multiple and Z is not
```
▶ curl --location --request GET 'localhost:8080/api/v1/riddle?x=3&y=6&z=4'
{
  "description": ""there is no solution to measure 4 with jugs with 3 and 6",
  "message": "invalid parameters"
}
```
