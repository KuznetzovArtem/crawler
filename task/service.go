package task

import (
	"crawler/domain"
	"crawler/repositories/logger"
	"errors"
	"fmt"
)

type ActionRunner interface {
	Do(request domain.Request) (*domain.Response, error)
}

// Service represents a Huawei service able to handle all the Huawei business logic.
type TaskService struct {
	log     logger.Logger
	actions map[string]ActionRunner
}

// SendApplication - controller for event when we need send request to PreBid server.
func (s *TaskService) Do(request domain.Request) (*domain.Response, error) {
	service, exist := s.actions[request.Action]
	if !exist {
		go s.log.Log(logger.Error, fmt.Sprintf("action doesn't support: %v", request.Action))
		return nil, errors.New("asd")
	}
	response, err := service.Do(request)
	return response, err
}

func NewService(
	log logger.Logger,
	actions map[string]ActionRunner,
) *TaskService {
	return &TaskService{
		log:     log,
		actions: actions,
	}
}
