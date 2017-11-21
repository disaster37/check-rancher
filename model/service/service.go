package modelService

import (
	"github.com/pkg/errors"
	rancherClient "github.com/rancher/go-rancher/v2"
	log "github.com/sirupsen/logrus"
)

// Rancher service data
type Service struct {
	service   *rancherClient.Service
	instances []rancherClient.Instance
}

type Services []Service

func NewService(service *rancherClient.Service) (*Service, error) {

	if service == nil {
		return nil, errors.New("Service can't be null")
	}

	log.Debug("Service: ", service.Name)

	newService := &Service{
		service:   service,
		instances: make([]rancherClient.Instance, 0),
	}

	return newService, nil
}

// Permit to set service
func (s *Service) SetService(service *rancherClient.Service) error {

	if service == nil {
		return errors.New("Service can't be null")
	}

	log.Debug("Service: ", service.Name)

	s.service = service

	return nil
}

// Permit to get service
func (s *Service) Service() *rancherClient.Service {

	return s.service
}

// Permit to set array of instance
func (s *Service) SetInstances(instances []rancherClient.Instance) error {

	if instances == nil {
		return errors.New("Instances can't be null")
	}

	log.Debug("Set %d instances", len(instances))

	s.instances = instances

	return nil
}

// Permit to add instance
func (s *Service) AddInstance(instance *rancherClient.Instance) error {

	if instance == nil {
		return errors.New("Instance can't be null")
	}

	log.Debug("Instance: ", instance.Name)

	s.instances = append(s.instances, *instance)

	return nil
}

// Permit to get all instances
func (s *Service) Instances() []rancherClient.Instance {

	return s.instances
}

// Permit to get instance on given index
func (s *Service) Instance(index int) (*rancherClient.Instance, error) {
	if index >= len(s.instances) {
		return nil, errors.New("Index is outbound of array")
	}

	return &s.Instances()[index], nil
}
