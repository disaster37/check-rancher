package rancherService

import (
	rancherClient "github.com/rancher/go-rancher/v2"
	log "github.com/sirupsen/logrus"
)

// Permit to retrive hosts linked to environment
func (r *Client) FindHostsByProjectLink(project *rancherClient.Project) ([]rancherClient.Host, error) {
	log.Debugf("Looking up Rancher hosts for project: %s", project.Name)
	hostCollection := rancherClient.HostCollection{}
	err := r.client.GetLink(project.Resource, "hosts", &hostCollection)
	if err != nil {
		return nil, err
	}

	if len(hostCollection.Data) == 0 {
		return nil, nil
	}

	log.Debugf("Found %d hosts link with project %s", len(hostCollection.Data), project.Name)

	return hostCollection.Data, nil
}
