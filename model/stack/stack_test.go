package modelStack

import (
	"github.com/disaster37/check-rancher/model/service"
	rancherClient "github.com/rancher/go-rancher/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Test the constructor
func TestNewStack(t *testing.T) {

	// Normal use case
	rs := &rancherClient.Stack{}
	s, _ := NewStack(rs)
	assert.Equal(t, rs, s.stack)
	assert.Equal(t, make(modelService.Services, 0), s.services)

	// Bad use case
	s, err := NewStack(nil)
	assert.Error(t, err)
}

// Test get and set stack
func TestGetSetStack(t *testing.T) {

	// Normal use case
	rs := &rancherClient.Stack{}
	s, _ := NewStack(&rancherClient.Stack{})
	s.SetStack(rs)
	assert.Equal(t, rs, s.Stack())

	// Bad use case
	err := s.SetStack(nil)
	assert.Error(t, err)
}

// Test get and set and add service
func TestGetSetAddService(t *testing.T) {
	s1, _ := modelService.NewService(&rancherClient.Service{})
	s2, _ := modelService.NewService(&rancherClient.Service{})

	listServices := modelService.Services{*s1}
	s, _ := NewStack(&rancherClient.Stack{})
	s.SetServices(listServices)
	assert.Equal(t, listServices, s.Services())

	s.AddService(s2)
	assert.Equal(t, 2, len(s.Services()))
	tempService, _ := s.Service(1)
	assert.Equal(t, s2, tempService)
}
