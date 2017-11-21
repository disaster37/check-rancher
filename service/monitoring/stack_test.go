package monitoringService

import (
	"github.com/disaster37/check-rancher/model/service"
	"github.com/disaster37/check-rancher/model/stack"
	rancherClient "github.com/rancher/go-rancher/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckStack(t *testing.T) {

	// When stack is null
	_, err := CheckStack(nil)
	assert.Error(t, err)

	// When stack is disabled
	rancherStack := &rancherClient.Stack{
		Name:  "test",
		State: "inactive",
	}
	stack, _ := modelStack.NewStack(rancherStack)
	monitoringData, _ := CheckStack(stack)
	assert.Equal(t, 0, monitoringData.Status())
	assert.Equal(t, 0, len(monitoringData.Perfdatas()))

	// When stack contains no services
	rancherStack = &rancherClient.Stack{
		Name:  "test",
		State: "active",
	}
	stack, _ = modelStack.NewStack(rancherStack)
	monitoringData, _ = CheckStack(stack)
	assert.Equal(t, 1, monitoringData.Status())
	assert.Equal(t, 2, len(monitoringData.Perfdatas()))

	// When stack contains all services disabled
	rancherStack = &rancherClient.Stack{
		Name:  "test",
		State: "active",
	}
	stack, _ = modelStack.NewStack(rancherStack)
	serviceRancher := &rancherClient.Service{
		Name:  "test",
		State: "inactive",
	}
	service, _ := modelService.NewService(serviceRancher)
	stack.AddService(service)
	monitoringData, _ = CheckStack(stack)
	assert.Equal(t, 0, monitoringData.Status())
	assert.Equal(t, 4, len(monitoringData.Perfdatas()))

	// When stack contains service OK
	rancherStack = &rancherClient.Stack{
		Name:  "test",
		State: "active",
	}
	stack, _ = modelStack.NewStack(rancherStack)
	serviceRancher = &rancherClient.Service{
		Name:        "test",
		State:       "active",
		HealthState: "healthy",
	}
	service, _ = modelService.NewService(serviceRancher)
	stack.AddService(service)
	monitoringData, _ = CheckStack(stack)
	assert.Equal(t, 0, monitoringData.Status())
	assert.Equal(t, 4, len(monitoringData.Perfdatas()))

	// When stack contains service KO with all instance Ko
	rancherStack = &rancherClient.Stack{
		Name:  "test",
		State: "active",
	}
	stack, _ = modelStack.NewStack(rancherStack)
	serviceRancher = &rancherClient.Service{
		Name:        "test",
		State:       "active",
		HealthState: "unhealthy",
	}
	service, _ = modelService.NewService(serviceRancher)
	instanceRancher := &rancherClient.Instance{
		Name:          "test",
		State:         "stopped",
		Transitioning: "no",
	}
	service.AddInstance(instanceRancher)
	stack.AddService(service)
	monitoringData, _ = CheckStack(stack)
	assert.Equal(t, 2, monitoringData.Status())
	assert.Equal(t, 4, len(monitoringData.Perfdatas()))

	// When stack contains service KO with one instance OK and other KO
	rancherStack = &rancherClient.Stack{
		Name:  "test",
		State: "active",
	}
	stack, _ = modelStack.NewStack(rancherStack)
	serviceRancher = &rancherClient.Service{
		Name:        "test",
		State:       "active",
		HealthState: "unhealthy",
	}
	service, _ = modelService.NewService(serviceRancher)
	instanceRancher1 := &rancherClient.Instance{
		Name:          "test",
		State:         "stopped",
		Transitioning: "no",
	}
	service.AddInstance(instanceRancher1)
	instanceRancher2 := &rancherClient.Instance{
		Name:          "test",
		State:         "running",
		Transitioning: "no",
	}
	service.AddInstance(instanceRancher2)
	stack.AddService(service)
	monitoringData, _ = CheckStack(stack)
	assert.Equal(t, 1, monitoringData.Status())
	assert.Equal(t, 4, len(monitoringData.Perfdatas()))
}
