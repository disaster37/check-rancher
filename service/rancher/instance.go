package rancherService

import (
	rancherClient "github.com/rancher/go-rancher/v2"
	log "github.com/sirupsen/logrus"
)

// Permit to retrive instances linked to service
func (r *Client) FindInstancesByServiceLink(service *rancherClient.Service) ([]rancherClient.Instance, error) {
	log.Debugf("Looking up Rancher instances for services: %s", service.Name)
	instanceCollection := rancherClient.InstanceCollection{}
	err := r.client.GetLink(service.Resource, "instances", &instanceCollection)
	if err != nil {
		return nil, err
	}

	if len(instanceCollection.Data) == 0 {
		return nil, nil
	}

	log.Debugf("Found %d instances link with service %s", len(instanceCollection.Data), service.Name)

	return instanceCollection.Data, nil
}
