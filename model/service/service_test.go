package modelService

import (
	rancherClient "github.com/rancher/go-rancher/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Test the constructor
func TestNewService(t *testing.T) {

	// Normal use case
	rs := &rancherClient.Service{}
	s, _ := NewService(rs)
	assert.Equal(t, rs, s.service)
	assert.Equal(t, make([]rancherClient.Instance, 0), s.instances)

	// Bad use case
	s, err := NewService(nil)
	assert.Error(t, err)
}

// Test get and set service
func TestGetSetService(t *testing.T) {

	// Normal use case
	rs := &rancherClient.Service{}
	s, _ := NewService(&rancherClient.Service{})
	s.SetService(rs)
	assert.Equal(t, rs, s.Service())

	//Bad use case
	err := s.SetService(nil)
	assert.Error(t, err)
}

// Test get and set and add instance
func TestGetSetAddInstance(t *testing.T) {
	// Normal use case
	i := &rancherClient.Instance{}
	i2 := &rancherClient.Instance{}
	listInstances := []rancherClient.Instance{*i}
	s, _ := NewService(&rancherClient.Service{})
	s.SetInstances(listInstances)
	assert.Equal(t, listInstances, s.Instances())

	s.AddInstance(i2)
	assert.Equal(t, 2, len(s.Instances()))
	tempInstance, _ := s.Instance(1)
	assert.Equal(t, i2, tempInstance)

	// Bad use cas
	err := s.AddInstance(nil)
	assert.Error(t, err)
}
