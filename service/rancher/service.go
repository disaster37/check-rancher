package rancherService

import (
	rancherClient "github.com/rancher/go-rancher/v2"
	log "github.com/sirupsen/logrus"
)

// Permit to retrive services linked to stack
func (r *Client) FindServicesByStackLink(stack *rancherClient.Stack) ([]rancherClient.Service, error) {
	log.Debugf("Looking up Rancher services for stack: %s", stack.Name)
	serviceCollection := rancherClient.ServiceCollection{}
	err := r.client.GetLink(stack.Resource, "services", &serviceCollection)
	if err != nil {
		return nil, err
	}

	if len(serviceCollection.Data) == 0 {
		return nil, nil
	}

	log.Debugf("Found %d services link with stack %s", len(serviceCollection.Data), stack.Name)

	return serviceCollection.Data, nil
}
