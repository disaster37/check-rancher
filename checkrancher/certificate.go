package checkrancher

import (
	"strings"
	"time"

	"github.com/disaster37/go-nagios"
	"github.com/pkg/errors"
	rancher "github.com/rancher/go-rancher/v2"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// CheckCertificates wrap cli with monitoring check
func CheckCertificates(c *cli.Context) error {

	monitorRancher, err := manageRancherGlobalParameters(c)
	if err != nil {
		return err
	}

	monitoringData, err := monitorRancher.CheckCertificates(c.String("project-name"), c.Int("warning-days"))
	if err != nil {
		return err
	}
	monitoringData.ToSdtOut()

	return nil

}

// CheckCertificates permit to check certificat validity
func (h *CheckRancher) CheckCertificates(projectName string, thresholdWarningDay int) (*nagiosPlugin.Monitoring, error) {

	if projectName == "" {
		return nil, errors.New("ProjectName can't be empty")
	}
	log.Debugf("ProjectName: %s", projectName)
	log.Debugf("ThresholdWarningDay: %d", thresholdWarningDay)

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

	// If project is disable
	if project.State == "inactive" {
		monitoringData.SetStatus(nagiosPlugin.STATUS_OK)
		monitoringData.AddMessage("Project %s is inactive", project.Name)
		monitoringData.AddPerfdata("nbCertificates", 0, "")
		monitoringData.AddPerfdata("nbCertificatesExpired", 0, "")
		return monitoringData, nil
	}

	// Get certificats link on project
	certificates := &rancher.CertificateCollection{}
	err = h.client.GetLink(project.Resource, "certificates", certificates)
	if err != nil {
		return nil, err
	}
	log.Debugf("Found %d Rancher certificates", len(certificates.Data))

	// No certificat on project
	if len(certificates.Data) == 0 {
		monitoringData.SetStatus(nagiosPlugin.STATUS_OK)
		monitoringData.AddMessage("All certificates are valid: %d/%d", len(certificates.Data), len(certificates.Data))
		monitoringData.AddPerfdata("nbCertificates", 0, "")
		monitoringData.AddPerfdata("nbCertificatesExpired", 0, "")
		return monitoringData, nil
	}

	// Check certificates
	currentDatetime := time.Now()
	dateThreshold := currentDatetime.AddDate(0, 0, thresholdWarningDay)
	dateLayout := "Mon Jan 02 15:04:05 MST 2006"
	failedCertificates := make([]rancher.Certificate, 0)

	for _, certificate := range certificates.Data {
		if certificate.State == "active" {
			certificateExpireAt, err := time.Parse(dateLayout, certificate.ExpiresAt)
			if err != nil {
				return nil, err
			}
			log.Debugf("CertificateExpireAt: %s (%s)", certificateExpireAt, certificate.Name)
			diffCritical := certificateExpireAt.Sub(currentDatetime)
			diffWarning := certificateExpireAt.Sub(dateThreshold)
			log.Debugf("diffCritical %f days for certificate %s", diffCritical.Hours()/24, certificate.Name)
			log.Debugf("diffWarning %f days for certificate %s", diffWarning.Hours()/24, certificate.Name)
			if diffCritical <= 0 {
				monitoringData.SetStatus(nagiosPlugin.STATUS_CRITICAL)
				failedCertificates = append(failedCertificates, certificate)
			} else if diffWarning <= 0 {
				monitoringData.SetStatus(nagiosPlugin.STATUS_WARNING)
				failedCertificates = append(failedCertificates, certificate)
			}
		} else {
			log.Debugf("Certificate %s is disable", certificate.Name)
		}
	}

	if len(failedCertificates) == 0 {
		monitoringData.SetStatus(nagiosPlugin.STATUS_OK)
		monitoringData.AddMessage("All certificates are valid: %d/%d", len(certificates.Data), len(certificates.Data))
	} else {
		monitoringData.AddMessage("Some certificates issue: %d/%d", len(failedCertificates), len(certificates.Data))
		for _, certificate := range failedCertificates {
			monitoringData.AddMessage("Certificate %s expire at %s (%s)", certificate.Name, certificate.ExpiresAt, strings.Join(certificate.SubjectAlternativeNames, ", "))
		}
	}

	monitoringData.AddPerfdata("nbCertificates", len(certificates.Data), "")
	monitoringData.AddPerfdata("nbCertificatesExpired", len(failedCertificates), "")

	return monitoringData, nil
}
