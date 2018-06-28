package monitoringService

import (
	"fmt"
	"github.com/disaster37/check-rancher/model/monitoring"
	"github.com/disaster37/check-rancher/model/stack"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Permit to check stack state from Rancher Stack data
func CheckStack(stack *modelStack.Stack) (*modelMonitoring.Monitoring, error) {

	if stack == nil {
		return nil, errors.New("Stack can't be null")
	}

	log.Debugf("Check state for stack %s", stack.Stack().Name)

	nbServiceOnline := 0

	monitoringData := modelMonitoring.NewMonitoring()

	stackName := stack.Stack().Name

	// Check if stack is stopped
	if stack.Stack().State == "active" {

		// Get perfdata
		nbService := len(stack.Services())
		monitoringData.AddPerfdata("nbService", nbService, "")
		nbInstance := 0
		for _, service := range stack.Services() {
			log.Debug("Service: ", service.Service().Name)
			nbInstance += len(service.Instances())
		}
		monitoringData.AddPerfdata("nbInstance", nbInstance, "")

		// Check if stack have services
		if len(stack.Services()) == 0 {
			monitoringData.SetStatus(modelMonitoring.STATUS_WARNING)
			monitoringData.AddMessage(fmt.Sprintf("Stack '%s' have no service", stackName))
			return monitoringData, nil
		}

		// Check services states
		for _, service := range stack.Services() {

			serviceName := service.Service().Name

			// Count the nbInstance online for the current service
			nbInstanceOnline := 0
			for _, instance := range service.Instances() {
				log.Debugf("Instance %s: %s", instance.Name, instance.State)
				if instance.State == "running" && instance.Transitioning == "no" {
					nbInstanceOnline++
				}
			}
			monitoringData.AddPerfdata(fmt.Sprintf("%s-nbInstanceOnline", serviceName), nbInstanceOnline, "")

			// Service not actif
			if service.Service().State == "inactive" || service.Service().State == "upgrading" {
				monitoringData.AddMessage(fmt.Sprintf("Service '%s' is %s (%d/%d instance online)", serviceName, service.Service().State, nbInstanceOnline, len(service.Instances())))
			} else {
				// Service actif
				if service.Service().HealthState != "healthy" {
					if service.Service().HealthState == "initializing" {
						monitoringData.AddMessage(fmt.Sprintf("Service '%s' is in initializing state (%d/%d instance online)", serviceName, nbInstanceOnline, len(service.Instances())))
					} else {
						monitoringData.AddMessage(fmt.Sprintf("Service '%s' has problem (%d/%d instance online)", serviceName, nbInstanceOnline, len(service.Instances())))
					}

					if nbInstanceOnline == 0 {
						monitoringData.SetStatus(modelMonitoring.STATUS_CRITICAL)
					} else {
						monitoringData.SetStatus(modelMonitoring.STATUS_WARNING)
						nbServiceOnline++
					}
				} else {
					nbServiceOnline++
				}
			}
		}
		monitoringData.AddPerfdata("nbServiceOnline", nbServiceOnline, "")

		if monitoringData.Status() == modelMonitoring.STATUS_OK {
			monitoringData.AddMessage(fmt.Sprintf("Stack '%s' work fine (%d services, %d instances)", stackName, nbService, nbInstance))
		}

	} else {
		monitoringData.AddMessage(fmt.Sprintf("Stack '%s' is not active. Current status is %s", stackName, stack.Stack().State))
		return monitoringData, nil
	}

	return monitoringData, nil

}
