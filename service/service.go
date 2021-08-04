package service

type AppError struct {
	Error   error
	Message string
	Code    int
}

// Service describes service to deal with devices.
type Service interface {
	// Health: returns server status
	Health() *HealthResponse
	// Riddle: Solves Water Jug Riddle
	Riddle(x, y, z int) (*RiddleResponse, *AppError)
}

type service struct {
	// if there were a repository, it would be stored here
}

// NewService creates new instance for devices service.
func NewService() *service {
	return &service{}
}