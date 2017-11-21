package modelProject

import (
	rancherClient "github.com/rancher/go-rancher/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Test the constructor
func TestNewProject(t *testing.T) {

	// Normal use case
	rp := &rancherClient.Project{}
	p, _ := NewProject(rp)
	assert.Equal(t, rp, p.project)
	assert.Equal(t, make([]rancherClient.Host, 0), p.hosts)

	// Bad use case
	p, err := NewProject(nil)
	assert.Error(t, err)
}

// Test get and set project
func TestGetSetProject(t *testing.T) {

	// Normal use case
	rp := &rancherClient.Project{}
	p, _ := NewProject(&rancherClient.Project{})
	p.SetProject(rp)
	assert.Equal(t, rp, p.Project())

	// Bad use case
	err := p.SetProject(nil)
	assert.Error(t, err)
}

// Test get and set and add host
func TestGetSetAddHost(t *testing.T) {
	h := &rancherClient.Host{}
	h2 := &rancherClient.Host{}
	listHosts := []rancherClient.Host{*h}
	p, _ := NewProject(&rancherClient.Project{})
	p.SetHosts(listHosts)
	assert.Equal(t, listHosts, p.Hosts())

	p.AddHost(h2)
	assert.Equal(t, 2, len(p.Hosts()))
	tempHost, _ := p.Host(1)
	assert.Equal(t, h2, tempHost)
}
