package rancherService

import (
	rancherClient "github.com/rancher/go-rancher/v2"
	log "github.com/sirupsen/logrus"
)

// Permit to retrieve an existing project by name
func (r *Client) FindProjectByName(name string) (*rancherClient.Project, error) {
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
