package checkrancher

import (
	"github.com/disaster37/go-nagios"
	"github.com/pkg/errors"
	rancher "github.com/rancher/go-rancher/v2"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// CheckHosts wrap cli with monitoring check
func CheckHosts(c *cli.Context) error {

	monitorRancher, err := manageRancherGlobalParameters(c)
	if err != nil {
		return err
	}

	monitoringData, err := monitorRancher.CheckHosts(c.String("project-name"))
	if err != nil {
		return err
	}
	monitoringData.ToSdtOut()

	return nil

}

// CheckHosts permit to check hosts state on project
func (h *CheckRancher) CheckHosts(projectName string) (*nagiosPlugin.Monitoring, error) {

	if projectName == "" {
		return nil, errors.New("ProjectName can't be empty")
	}
	log.Debugf("ProjectName: %s", projectName)

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
		monitoringData.AddMessage("Project %s is inactive", project.Name)
		monitoringData.AddPerfdata("nbHosts", 0, "")
		monitoringData.AddPerfdata("nbFailedHosts", 0, "")
		monitoringData.AddPerfdata("nbInactiveHosts", 0, "")
		return monitoringData, nil
	}

	// Get hosts
	hosts := &rancher.HostCollection{}
	err = h.client.GetLink(project.Resource, "hosts", hosts)
	if err != nil {
		return nil, err
	}

	// When no host
	if len(hosts.Data) == 0 {
		monitoringData.SetStatus(nagiosPlugin.STATUS_OK)
		monitoringData.AddMessage("All hosts are active: %d/%d", 0, 0)
		monitoringData.AddPerfdata("nbHosts", 0, "")
		monitoringData.AddPerfdata("nbFailedHosts", 0, "")
		monitoringData.AddPerfdata("nbInactiveHosts", 0, "")
		return monitoringData, nil
	}

	// Check host state
	inactiveHosts := make([]rancher.Host, 0)
	failedHosts := make([]rancher.Host, 0)
	for _, host := range hosts.Data {
		if host.State == "inactive" {
			inactiveHosts = append(inactiveHosts, host)
		} else if host.State != "active" {
			failedHosts = append(failedHosts, host)
		}
	}

	if len(failedHosts) > 0 {
		monitoringData.SetStatus(nagiosPlugin.STATUS_CRITICAL)
		monitoringData.AddMessage("Some hosts failed: %d/%d", len(failedHosts), len(hosts.Data))
		for _, host := range failedHosts {
			monitoringData.AddMessage("Host %s in state %s", host.Name, host.State)
		}
	} else {
		monitoringData.AddMessage("All hosts are active: %d/%d", len(hosts.Data), len(hosts.Data))
	}
	if len(inactiveHosts) > 0 {
		monitoringData.AddMessage("Some hosts are inactive: %d/%d", len(inactiveHosts), len(hosts.Data))
		for _, host := range inactiveHosts {
			monitoringData.AddMessage("Host %s is inactive", host.Name)
		}
	}

	// Compute perfdata
	monitoringData.AddPerfdata("nbHosts", len(hosts.Data), "")
	monitoringData.AddPerfdata("nbFailedHosts", len(failedHosts), "")
	monitoringData.AddPerfdata("nbInactiveHosts", len(inactiveHosts), "")

	return monitoringData, nil
}
