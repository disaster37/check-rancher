package modelStack

import (
	"github.com/disaster37/check-rancher/model/service"
	"github.com/pkg/errors"
	rancherClient "github.com/rancher/go-rancher/v2"
	log "github.com/sirupsen/logrus"
)

// Rancher stack data
type Stack struct {
	stack    *rancherClient.Stack
	services modelService.Services
}

// Construct stack
func NewStack(stack *rancherClient.Stack) (*Stack, error) {

	if stack == nil {
		return nil, errors.New("Stack can't be null")
	}

	log.Debugf("Stack: %s", stack.Name)

	newStack := &Stack{
		stack:    stack,
		services: make(modelService.Services, 0, 3),
	}

	return newStack, nil

}

// Get stack
func (s *Stack) Stack() *rancherClient.Stack {
	return s.stack
}

// Set stack
func (s *Stack) SetStack(stack *rancherClient.Stack) error {

	if stack == nil {
		return errors.New("Stack can't be null")
	}

	log.Debugf("Set stack %s", stack.Name)

	s.stack = stack

	return nil
}

// Set list of service
func (s *Stack) SetServices(services modelService.Services) error {

	if services == nil {
		return errors.New("Services can't be null")
	}

	log.Debugf("Set %s services", len(services))

	s.services = services

	return nil
}

// Add service
func (s *Stack) AddService(service *modelService.Service) error {

	if service == nil {
		return errors.New("Service can't be null")
	}

	log.Debugf("Add service %s", service.Service().Name)

	s.services = append(s.services, *service)

	return nil
}

// Get list of service
func (s *Stack) Services() modelService.Services {
	return s.services
}

// Get one service by index
func (s *Stack) Service(index int) (*modelService.Service, error) {

	if index >= len(s.services) {
		return nil, errors.New("Index is outbound of array")
	}

	return &s.Services()[index], nil
}
