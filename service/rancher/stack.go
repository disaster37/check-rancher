package rancherService

import (
	"fmt"
	"github.com/disaster37/check-rancher/model/service"
	"github.com/disaster37/check-rancher/model/stack"
	"github.com/pkg/errors"
	rancherClient "github.com/rancher/go-rancher/v2"
	log "github.com/sirupsen/logrus"
)

// Permit to retrieve an existing stack by name
func (r *Client) FindStackByNameOnProject(name string, project *rancherClient.Project) (*rancherClient.Stack, error) {
	log.Debugf("Looking up Rancher stack by name %s on project %s", name, project.Name)

	stacks, err := r.client.Stack.List(&rancherClient.ListOpts{
		Filters: map[string]interface{}{
			"name":         name,
			"accountId":    project.Id,
			"removed_null": nil,
		},
	})

	if err != nil {
		return nil, err
	}

	if len(stacks.Data) == 0 {
		return nil, nil
	}

	log.Debugf("Found existing Rancher Stack by name: %s", name)
	return &stacks.Data[0], nil
}

// Permit to load Rancher stack by it's name
func (r *Client) LoadStackByNameOnProject(name string, project *rancherClient.Project) (*modelStack.Stack, error) {

	log.Debugf("Load Rancher stack by name: %s", name)

	// Get rancher stack
	rancherStack, err := r.FindStackByNameOnProject(name, project)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error appear when try to get Rancher stack %s", name))
	}

	// Stack not found
	if rancherStack == nil {
		return nil, errors.New(fmt.Sprintf("Stack %s not found", name))
	}

	stack, err := modelStack.NewStack(rancherStack)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error appear when try to create new stack %s", name))
	}

	// Get services linked to stack
	services, err := r.FindServicesByStackLink(rancherStack)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error apear when try to get Rancher services from stack %s", name))
	}

	if services != nil {
		// Get instances for each service
		for idx, rancherService := range services {
			log.Debug("Service: ", rancherService.Name)

			if rancherService.Type == "service" {

				service, err := modelService.NewService(&services[idx])
				if err != nil {
					return nil, errors.Wrap(err, fmt.Sprintf("Error appear when try to create new service on stack %s", name))
				}

				instances, err := r.FindInstancesByServiceLink(&services[idx])
				if err != nil {
					return nil, errors.Wrap(err, fmt.Sprintf("Error apear when try to get Rancher instance from service %s", rancherService.Name))
				}

				service.SetInstances(instances)

				stack.AddService(service)

			}

		}
	}

	return stack, nil

}
