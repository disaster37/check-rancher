package monitoringService

import (
	rancherClient "github.com/rancher/go-rancher/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckCertificates(t *testing.T) {

	// When list of certificates is null
	_, err := CheckCertificates(nil, 0)
	assert.Error(t, err)

	// When list of certificates is empty
	certificates := make([]rancherClient.Certificate, 0)
	monitoringData, _ := CheckCertificates(certificates, 0)
	assert.Equal(t, 0, monitoringData.Status())
	assert.Equal(t, 2, len(monitoringData.Perfdatas()))
	
	// When all certificates are disabled
	certificates = make([]rancherClient.Certificate, 0)
	certificate := rancherClient.Certificate{
		Name:      "test",
		State:     "inactive",
	}
	certificates = append(certificates, certificate)
	monitoringData, _ = CheckCertificates(certificates, 0)
	assert.Equal(t, 0, monitoringData.Status())
	assert.Equal(t, 2, len(monitoringData.Perfdatas()))

	// When certificate not yet in threshold
	certificates = make([]rancherClient.Certificate, 0)
	certificate = rancherClient.Certificate{
		Name:      "test",
		State:     "active",
		ExpiresAt: "Mon Feb 12 11:13:07 UTC 2118",
	}
	certificates = append(certificates, certificate)
	monitoringData, _ = CheckCertificates(certificates, 0)
	assert.Equal(t, 0, monitoringData.Status())
	assert.Equal(t, 2, len(monitoringData.Perfdatas()))

	// When certificates in threshold
	certificates = make([]rancherClient.Certificate, 0)
	certificate = rancherClient.Certificate{
		Name:      "test",
		State:     "active",
		ExpiresAt: "Mon Feb 12 11:13:07 UTC 2118",
	}
	certificates = append(certificates, certificate)
	monitoringData, _ = CheckCertificates(certificates, 365000)
	assert.Equal(t, 1, monitoringData.Status())
	assert.Equal(t, 2, len(monitoringData.Perfdatas()))

	// When certificate is expired
	certificates = make([]rancherClient.Certificate, 0)
	certificate = rancherClient.Certificate{
		Name:      "test",
		State:     "active",
		ExpiresAt: "Mon Feb 12 11:13:07 UTC 2017",
	}
	certificates = append(certificates, certificate)
	monitoringData, _ = CheckCertificates(certificates, 365)
	assert.Equal(t, 2, monitoringData.Status())
	assert.Equal(t, 2, len(monitoringData.Perfdatas()))

	// check the priority status
	certificates = make([]rancherClient.Certificate, 0)
	certificate1 := rancherClient.Certificate{
		Name:      "test",
		State:     "active",
		ExpiresAt: "Mon Feb 12 11:13:07 UTC 2017",
	}
	certificate2 := rancherClient.Certificate{
		Name:      "test",
		State:     "active",
		ExpiresAt: "Mon Feb 12 11:13:07 UTC 2118",
	}
	certificates = append(certificates, certificate1)
	certificates = append(certificates, certificate2)
	monitoringData, _ = CheckCertificates(certificates, 365000)
	assert.Equal(t, 2, monitoringData.Status())
	assert.Equal(t, 2, len(monitoringData.Perfdatas()))

}
