package monitoringService

import (
	"github.com/disaster37/check-rancher/model/monitoring"
	"github.com/disaster37/check-rancher/model/project"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Permit to check the project's hosts state from Rancher project data
func CheckHostsProject(project *modelProject.Project) (*modelMonitoring.Monitoring, error) {

	if project == nil {
		return nil, errors.New("Project can't be null")
	}

	log.Debugf("Check state for hosts's project %s", project.Project().Name)

	monitoringData := modelMonitoring.NewMonitoring()
	projectName := project.Project().Name

	// Check if project is enable
	if project.Project().State == "active" {

		// Get perfdata
		nbHost := len(project.Hosts())
		monitoringData.AddPerfdata("nbHost", nbHost, "")

		// Check if project have hosts
		if len(project.Hosts()) == 0 {
			monitoringData.SetStatus(modelMonitoring.STATUS_WARNING)
			monitoringData.AddMessage(fmt.Sprintf("Project '%s' have no host", projectName))
			// Add perfdata to follow the host online
			monitoringData.AddPerfdata("nbHostOnline", 0, "")
			return monitoringData, nil
		}

		// We check the hosts state
		nbHostOnline := 0
		for _, host := range project.Hosts() {
			if host.State == "inactive" {
				monitoringData.AddMessage(fmt.Sprintf("Host '%s' is disabled", host.Name))
			} else if host.State != "active" {
				monitoringData.AddMessage(fmt.Sprintf("Host '%s' have problem", host.Name))
				monitoringData.SetStatus(modelMonitoring.STATUS_CRITICAL)
			} else {
				nbHostOnline++
			}
		}

		// Add perfdata to follow the host online
		monitoringData.AddPerfdata("nbHostOnline", nbHostOnline, "")

		// All host work fine
		if nbHostOnline == len(project.Hosts()) {
			monitoringData.AddMessage(fmt.Sprintf("All hosts work fine (%d/%d hosts online)", nbHostOnline, len(project.Hosts())))
		} else {
			if nbHostOnline == 0 {
				// Maybee all host are disabled
				monitoringData.SetStatus(modelMonitoring.STATUS_WARNING)
			}
			monitoringData.AddMessage(fmt.Sprintf("Only %d/%d hosts online", nbHostOnline, len(project.Hosts())))
		}

	} else {
		monitoringData.AddMessage(fmt.Sprintf("Project '%s' is not active. Current status is %s", projectName, project.Project().State))
		return monitoringData, nil
	}

	return monitoringData, nil

}
