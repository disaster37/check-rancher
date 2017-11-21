package rancherService

import (
	"github.com/disaster37/check-rancher/model/project"
	"fmt"
	"github.com/pkg/errors"
	rancherClient "github.com/rancher/go-rancher/v2"
	log "github.com/sirupsen/logrus"
)

// Permit to retrieve an existing project by name
func (r *Client) FindProjecttByName(name string) (*rancherClient.Project, error) {
	log.Debugf("Looking up Rancher project by name: %s", name)

	projects, err := r.client.Project.List(&rancherClient.ListOpts{
		Filters: map[string]interface{}{
			"name":         name,
			"removed_null": nil,
		},
	})

	if err != nil {
		return nil, err
	}

	if len(projects.Data) == 0 {
		return nil, nil
	}

	log.Debugf("Found existing Rancher project by name: %s", name)
	return &projects.Data[0], nil
}

// Permit to load project and hosts linked to it
func (r *Client) LoadProjectByName(name string) (*modelProject.Project, error) {

	log.Debugf("Load Rancher project by name: %s", name)

	// Get rancher project
	rancherProject, err := r.FindProjecttByName(name)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error appear when try to get Rancher project %s", name))
	}

	if rancherProject == nil {
		return nil, errors.New(fmt.Sprintf("Project %s not found", name))
	}

	project, err := modelProject.NewProject(rancherProject)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error appear when try to create new project %s", name))
	}

	// Get hosts linked to project
	hosts, err := r.FindHostsByProjectLink(rancherProject)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error apear when try to get Rancher hosts from project %s", name))
	}

	project.SetHosts(hosts)

	return project, nil

}
