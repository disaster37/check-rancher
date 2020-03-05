package checkrancher

import (
	"fmt"

	"github.com/disaster37/go-nagios"
	"github.com/pkg/errors"
	rancher "github.com/rancher/go-rancher/v2"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// CheckStack wrap cli with monitoring check
func CheckStack(c *cli.Context) error {

	monitorRancher, err := manageRancherGlobalParameters(c)
	if err != nil {
		return err
	}

	if c.String("stack-name") == "" {
		return errors.New("You must provide --stack-name parameter")
	}

	monitoringData, err := monitorRancher.CheckStack(c.String("project-name"), c.String("stack-name"))
	if err != nil {
		return err
	}
	monitoringData.ToSdtOut()

	return nil

}

// CheckStack permit to check stack state on project
func (h *CheckRancher) CheckStack(projectName string, stackName string) (*nagiosPlugin.Monitoring, error) {

	if projectName == "" {
		return nil, errors.New("ProjectName can't be empty")
	}
	if stackName == "" {
		return nil, errors.New("StackName can't be empty")
	}
	log.Debugf("ProjectName: %s", projectName)
	log.Debugf("StackName: %s", stackName)

	monitoringData := nagiosPlugin.NewMonitoring()

	// Get the project
	project, err := h.FindProjectByName(projectName)
	if err != nil {
		return nil, err
	}
	if project == nil {
		monitoringData.SetStatus(nagiosPlugin.STATUS_UNKNOWN)
		monitoringData.AddMessage("Project %s not found", projectName)
		return monitoringData, nil
	}

	// When project is disabled
	if project.State == "inactive" {
		monitoringData.SetStatus(nagiosPlugin.STATUS_OK)
		monitoringData.AddMessage("Project %s is inactive", project.Name)
		monitoringData.AddPerfdata("nbFailedServices", 0, "")
		monitoringData.AddPerfdata("nbUpgradedServices", 0, "")
		monitoringData.AddPerfdata("nbInactiveServices", 0, "")
		return monitoringData, nil
	}

	// Get stack
	stacks := &rancher.StackCollection{}
	err = h.client.GetLink(project.Resource, "stacks", stacks)
	if err != nil {
		return nil, err
	}

	var stack *rancher.Stack
	for _, stackTmp := range stacks.Data {
		if stackTmp.Name == stackName {
			stack = &stackTmp
			break
		}
	}

	if stack == nil {
		monitoringData.SetStatus(nagiosPlugin.STATUS_UNKNOWN)
		monitoringData.AddMessage("Stack %s not found", stackName)
		return monitoringData, nil
	}

	// When stack is disable
	if stack.State == "inactive" {
		monitoringData.SetStatus(nagiosPlugin.STATUS_OK)
		monitoringData.AddMessage("Stack %s is inactive", stack.Name)
		monitoringData.AddPerfdata("nbFailedServices", 0, "")
		monitoringData.AddPerfdata("nbUpgradedServices", 0, "")
		monitoringData.AddPerfdata("nbInactiveServices", 0, "")
		return monitoringData, nil
	}

	// Load services
	services := &rancher.ServiceCollection{}
	err = h.client.GetLink(stack.Resource, "services", services)
	if err != nil {
		return nil, err
	}

	// When no services
	if len(services.Data) == 0 {
		monitoringData.SetStatus(nagiosPlugin.STATUS_WARNING)
		monitoringData.AddMessage("No service found in stack %s", stack.Name)
		monitoringData.AddPerfdata("nbFailedServices", 0, "")
		monitoringData.AddPerfdata("nbUpgradedServices", 0, "")
		monitoringData.AddPerfdata("nbInactiveServices", 0, "")
		return monitoringData, nil
	}

	// Check service state
	disabledServices := make([]rancher.Service, 0)
	upgradedServices := make([]rancher.Service, 0)
	failedServices := make([]rancher.Service, 0)
	for _, service := range services.Data {
		switch service.State {
		case "inactive":
			disabledServices = append(disabledServices, service)
			instances := &rancher.InstanceCollection{}
			err = h.client.GetLink(service.Resource, "instances", instances)
			if err != nil {
				return nil, err
			}
			monitoringData.AddPerfdata(fmt.Sprintf("%s-nbScale", service.Name), int(service.CurrentScale), "")
			monitoringData.AddPerfdata(fmt.Sprintf("%s-nbInstances", service.Name), len(instances.Data), "")
			monitoringData.AddPerfdata(fmt.Sprintf("%s-nbFailedInstances", service.Name), 0, "")
		case "upgrading":
			upgradedServices = append(upgradedServices, service)
			instances := &rancher.InstanceCollection{}
			err = h.client.GetLink(service.Resource, "instances", instances)
			if err != nil {
				return nil, err
			}
			monitoringData.AddPerfdata(fmt.Sprintf("%s-nbScale", service.Name), int(service.CurrentScale), "")
			monitoringData.AddPerfdata(fmt.Sprintf("%s-nbInstances", service.Name), len(instances.Data), "")
			monitoringData.AddPerfdata(fmt.Sprintf("%s-nbFailedInstances", service.Name), 0, "")
		default:
			if service.HealthState != "healthy" {
				failedServices = append(failedServices, service)
			} else {
				instances := &rancher.InstanceCollection{}
				err = h.client.GetLink(service.Resource, "instances", instances)
				if err != nil {
					return nil, err
				}
				monitoringData.AddPerfdata(fmt.Sprintf("%s-nbScale", service.Name), int(service.CurrentScale), "")
				monitoringData.AddPerfdata(fmt.Sprintf("%s-nbInstances", service.Name), len(instances.Data), "")
				monitoringData.AddPerfdata(fmt.Sprintf("%s-nbFailedInstances", service.Name), 0, "")
			}
		}
	}

	if len(failedServices) > 0 {
		monitoringData.AddMessage("Some services failed: %d/%d", len(failedServices), len(services.Data))
		for _, service := range failedServices {
			// Load instances
			instances := &rancher.InstanceCollection{}
			err = h.client.GetLink(service.Resource, "instances", instances)
			if err != nil {
				return nil, err
			}

			if len(instances.Data) == 0 {
				monitoringData.SetStatus(nagiosPlugin.STATUS_CRITICAL)
				monitoringData.AddMessage("Service %s is critical: %d/%d instances", service.Name, 0, service.CurrentScale)
				monitoringData.AddPerfdata(fmt.Sprintf("%s-nbScale", service.Name), int(service.CurrentScale), "")
				monitoringData.AddPerfdata(fmt.Sprintf("%s-nbInstances", service.Name), 0, "")
				monitoringData.AddPerfdata(fmt.Sprintf("%s-nbFailedInstances", service.Name), int(service.CurrentScale), "")
			} else {
				nbOnlineInstance := 0
				for _, instance := range instances.Data {
					if instance.State == "running" && instance.Transitioning == "no" {
						nbOnlineInstance++
					}
				}
				monitoringData.AddPerfdata(fmt.Sprintf("%s-nbScale", service.Name), int(service.CurrentScale), "")
				monitoringData.AddPerfdata(fmt.Sprintf("%s-nbInstances", service.Name), len(instances.Data), "")
				monitoringData.AddPerfdata(fmt.Sprintf("%s-nbFailedInstances", service.Name), len(instances.Data)-nbOnlineInstance, "")

				if nbOnlineInstance == 0 {
					monitoringData.SetStatus(nagiosPlugin.STATUS_CRITICAL)
					monitoringData.AddMessage("Service %s is critical: %d/%d instances", service.Name, 0, len(instances.Data))
				} else {
					monitoringData.SetStatus(nagiosPlugin.STATUS_WARNING)
					monitoringData.AddMessage("Service %s is degraded: %d/%d instances", service.Name, nbOnlineInstance, len(instances.Data))
				}
			}
		}
	} else {
		monitoringData.SetStatus(nagiosPlugin.STATUS_OK)
		monitoringData.AddMessage("All services are ok: %d/%d", len(services.Data), len(services.Data))
	}

	if len(upgradedServices) > 0 {
		monitoringData.AddMessage("Some service are upgrading: %d/%d", len(upgradedServices), len(services.Data))
		for _, service := range upgradedServices {
			monitoringData.AddMessage("Service %s is upgrading", service.Name)
		}
	}

	if len(disabledServices) > 0 {
		monitoringData.AddMessage("Some service are inactive: %d/%d", len(disabledServices), len(services.Data))
		for _, service := range disabledServices {
			monitoringData.AddMessage("Service %s is inactive", service.Name)
		}
	}

	// Compute perfdata
	monitoringData.AddPerfdata("nbServices", len(services.Data), "")
	monitoringData.AddPerfdata("nbFailedServices", len(failedServices), "")
	monitoringData.AddPerfdata("nbUpgradedServices", len(upgradedServices), "")
	monitoringData.AddPerfdata("nbInactiveServices", len(disabledServices), "")

	return monitoringData, nil
}
