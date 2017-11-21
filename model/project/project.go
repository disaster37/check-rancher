package modelProject

import (
	"github.com/pkg/errors"
	rancherClient "github.com/rancher/go-rancher/v2"
	log "github.com/sirupsen/logrus"
)

// Rancher project data
type Project struct {
	project *rancherClient.Project
	hosts   []rancherClient.Host
}

// Constructor
func NewProject(project *rancherClient.Project) (*Project, error) {

	if project == nil {
		return nil, errors.New("Project can't be null")
	}
	log.Debug("Project: ", project.Name)

	p := &Project{
		project: project,
		hosts:   make([]rancherClient.Host, 0),
	}

	return p, nil
}

// Set project
func (p *Project) SetProject(project *rancherClient.Project) error {

	if project == nil {
		return errors.New("Project can't be null")
	}
	log.Debug("Project: ", project.Name)

	p.project = project

	return nil

}

// Get project
func (p *Project) Project() *rancherClient.Project {
	return p.project
}

// Set hosts
func (p *Project) SetHosts(hosts []rancherClient.Host) error {

	if hosts == nil {
		return errors.New("Hosts can't be null")
	}

	log.Debugf("Add %d hosts", len(hosts))

	p.hosts = hosts

	return nil
}

// Add host
func (p *Project) AddHost(host *rancherClient.Host) error {

	if host == nil {
		return errors.New("Host can't be null")
	}

	log.Debug("Host: ", host.Name)

	p.hosts = append(p.hosts, *host)

	return nil
}

// Get hosts
func (p *Project) Hosts() []rancherClient.Host {
	return p.hosts
}

// Get host
func (p *Project) Host(index int) (*rancherClient.Host, error) {
	if index >= len(p.hosts) {
		return nil, errors.New("Index is outbound of array")
	}

	return &p.Hosts()[index], nil
}
