package monitoringService

import (
	"fmt"
	"github.com/disaster37/check-rancher/model/monitoring"
	"github.com/pkg/errors"
	rancherClient "github.com/rancher/go-rancher/v2"
	log "github.com/sirupsen/logrus"
)

// Permit to check the hosts state
func CheckHosts(hosts []rancherClient.Host) (*modelMonitoring.Monitoring, error) {

	if hosts == nil {
		return nil, errors.New("Hosts can't be null")
	}

	log.Debugf("Check state for hosts")

	monitoringData := modelMonitoring.NewMonitoring()

	// Get perfdata
	nbHost := len(hosts)
	monitoringData.AddPerfdata("nbHost", nbHost, "")

	// Check if hosts is not empty
	if len(hosts) == 0 {
		monitoringData.SetStatus(modelMonitoring.STATUS_WARNING)
		monitoringData.AddMessage(fmt.Sprintf("No host"))
		// Add perfdata to follow the host online
		monitoringData.AddPerfdata("nbHostOnline", 0, "")
		return monitoringData, nil
	}

	// We check the hosts state
	nbHostOnline := 0
	for _, host := range hosts {
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
	if nbHostOnline == len(hosts) {
		monitoringData.AddMessage(fmt.Sprintf("All hosts work fine (%d/%d hosts online)", nbHostOnline, len(hosts)))
	} else {
		if nbHostOnline == 0 {
			// Maybee all host are disabled
			monitoringData.SetStatus(modelMonitoring.STATUS_WARNING)
		}
		monitoringData.AddMessage(fmt.Sprintf("Only %d/%d hosts online", nbHostOnline, len(hosts)))
	}

	return monitoringData, nil

}
