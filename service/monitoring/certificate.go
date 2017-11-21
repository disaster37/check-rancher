package monitoringService

import (
	"github.com/disaster37/check-rancher/model/monitoring"
	"fmt"
	"github.com/pkg/errors"
	rancherClient "github.com/rancher/go-rancher/v2"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

// Permit to check the project's hosts state from Rancher project data
func CheckCertificates(certificates []rancherClient.Certificate, warningThresholdInDays int) (*modelMonitoring.Monitoring, error) {

	if certificates == nil {
		return nil, errors.New("Certificates can't be null")
	}

	log.Debugf("Check state for certificates")

	monitoringData := modelMonitoring.NewMonitoring()

	currentDatetime := time.Now()
	dateThreshold := currentDatetime.AddDate(0, 0, warningThresholdInDays)
	dateLayout := "Mon Jan 02 15:04:05 MST 2006"
	nbCertificateActif := 0
	nbCertificateExpire := 0

	for _, certificate := range certificates {
		if certificate.State == "active" {
			nbCertificateActif++
			certificateExpireAt, err := time.Parse(dateLayout, certificate.ExpiresAt)
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("Error apear when try to parse date %s", certificate.ExpiresAt))
			}
			log.Debugf("CertificateExpireAt: %s (%s)", certificateExpireAt, certificate.Name)
			diffCritical := certificateExpireAt.Sub(currentDatetime)
			diffWarning := certificateExpireAt.Sub(dateThreshold)
			log.Debugf("diffCritical %d days for certificate %s", diffCritical.Hours()/24, certificate.Name)
			log.Debugf("diffWarning %d days for certificate %s", diffWarning.Hours()/24, certificate.Name)
			if diffCritical <= 0 {
				monitoringData.SetStatus(modelMonitoring.STATUS_CRITICAL)
				nbCertificateExpire++
			} else if diffWarning <= 0 {
				monitoringData.SetStatus(modelMonitoring.STATUS_WARNING)
			}

			if monitoringData.Status() != modelMonitoring.STATUS_OK {
				monitoringData.AddMessage(fmt.Sprintf("Certificate %s expire at %s (%s)", certificate.Name, certificate.ExpiresAt, strings.Join(certificate.SubjectAlternativeNames, ", ")))
			}

		} else {
			log.Debugf("Certificate %s is disable", certificate.Name)
		}
	}

	monitoringData.AddPerfdata("nbActifCertificate", nbCertificateActif, "")
	monitoringData.AddPerfdata("nbExpiredCertificate", nbCertificateExpire, "")

	if monitoringData.Status() == modelMonitoring.STATUS_OK {
		monitoringData.AddMessage(fmt.Sprintf("All certificates are validated (%d certficates actif)", nbCertificateActif))
	}

	return monitoringData, nil
}
