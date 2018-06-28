package monitoringService

import (
	rancherClient "github.com/rancher/go-rancher/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckHosts(t *testing.T) {

	// When hosts is null
	_, err := CheckHosts(nil)
	assert.Error(t, err)

	// When hosts list is empty
	hosts := make([]rancherClient.Host, 0)
	monitoringData, err := CheckHosts(hosts)
	assert.Equal(t, 1, monitoringData.Status())
	assert.Equal(t, 2, len(monitoringData.Perfdatas()))

	// When contains only host disabled
	hosts = make([]rancherClient.Host, 0)
	host := rancherClient.Host{
		Name:  "Test1",
		State: "inactive",
	}
	hosts = append(hosts, host)
	monitoringData, err = CheckHosts(hosts)
	assert.Equal(t, 1, monitoringData.Status())
	assert.Equal(t, 2, len(monitoringData.Perfdatas()))

	// When all host work fine
	hosts = make([]rancherClient.Host, 0)
	host = rancherClient.Host{
		Name:  "Test1",
		State: "active",
	}
	hosts = append(hosts, host)
	monitoringData, err = CheckHosts(hosts)
	assert.Equal(t, 0, monitoringData.Status())
	assert.Equal(t, 2, len(monitoringData.Perfdatas()))

	// When host has problem
	hosts = make([]rancherClient.Host, 0)
	host = rancherClient.Host{
		Name:  "Test1",
		State: "disconnected",
	}
	hosts = append(hosts, host)
	monitoringData, err = CheckHosts(hosts)
	assert.Equal(t, 2, monitoringData.Status())
	assert.Equal(t, 2, len(monitoringData.Perfdatas()))

	// When host work fine with host disabled
	hosts = make([]rancherClient.Host, 0)
	host1 := rancherClient.Host{
		Name:  "Test1",
		State: "inactive",
	}
	hosts = append(hosts, host1)
	host2 := rancherClient.Host{
		Name:  "Test1",
		State: "active",
	}
	hosts = append(hosts, host2)
	monitoringData, err = CheckHosts(hosts)
	assert.Equal(t, 0, monitoringData.Status())
	assert.Equal(t, 2, len(monitoringData.Perfdatas()))

	// When host has problem with host disabled
	hosts = make([]rancherClient.Host, 0)
	host1 = rancherClient.Host{
		Name:  "Test1",
		State: "inactive",
	}
	hosts = append(hosts, host1)
	host2 = rancherClient.Host{
		Name:  "Test1",
		State: "disconnected",
	}
	hosts = append(hosts, host2)
	monitoringData, err = CheckHosts(hosts)
	assert.Equal(t, 2, monitoringData.Status())
	assert.Equal(t, 2, len(monitoringData.Perfdatas()))

}
