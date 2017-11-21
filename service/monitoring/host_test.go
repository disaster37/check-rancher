package monitoringService

import (
	"github.com/disaster37/check-rancher/model/project"
	rancherClient "github.com/rancher/go-rancher/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckHostsProject(t *testing.T) {

	// When project is nul
	_, err := CheckHostsProject(nil)
	assert.Error(t, err)

	// When project doesn't contains host
	rancherProject := &rancherClient.Project{
		Name:  "Test",
		State: "active",
	}
	project, _ := modelProject.NewProject(rancherProject)
	monitoringData, err := CheckHostsProject(project)
	assert.Equal(t, 1, monitoringData.Status())
	assert.Equal(t, 2, len(monitoringData.Perfdatas()))

	// When project contains only host disabled
	rancherProject = &rancherClient.Project{
		Name:  "Test",
		State: "active",
	}
	project, _ = modelProject.NewProject(rancherProject)
	host := &rancherClient.Host{
		Name:  "Test1",
		State: "inactive",
	}
	project.AddHost(host)
	monitoringData, err = CheckHostsProject(project)
	assert.Equal(t, 1, monitoringData.Status())
	assert.Equal(t, 2, len(monitoringData.Perfdatas()))

	// When all host work fine
	rancherProject = &rancherClient.Project{
		Name:  "Test",
		State: "active",
	}
	project, _ = modelProject.NewProject(rancherProject)
	host = &rancherClient.Host{
		Name:  "Test1",
		State: "active",
	}
	project.AddHost(host)
	monitoringData, err = CheckHostsProject(project)
	assert.Equal(t, 0, monitoringData.Status())
	assert.Equal(t, 2, len(monitoringData.Perfdatas()))

	// When host has problem
	rancherProject = &rancherClient.Project{
		Name:  "Test",
		State: "active",
	}
	project, _ = modelProject.NewProject(rancherProject)
	host = &rancherClient.Host{
		Name:  "Test1",
		State: "disconnected",
	}
	project.AddHost(host)
	monitoringData, err = CheckHostsProject(project)
	assert.Equal(t, 2, monitoringData.Status())
	assert.Equal(t, 2, len(monitoringData.Perfdatas()))

	// When host work fine with host disabled
	rancherProject = &rancherClient.Project{
		Name:  "Test",
		State: "active",
	}
	project, _ = modelProject.NewProject(rancherProject)
	host1 := &rancherClient.Host{
		Name:  "Test1",
		State: "inactive",
	}
	project.AddHost(host1)
	host2 := &rancherClient.Host{
		Name:  "Test1",
		State: "active",
	}
	project.AddHost(host2)
	monitoringData, err = CheckHostsProject(project)
	assert.Equal(t, 0, monitoringData.Status())
	assert.Equal(t, 2, len(monitoringData.Perfdatas()))

	// When host has problem with host disabled
	rancherProject = &rancherClient.Project{
		Name:  "Test",
		State: "active",
	}
	project, _ = modelProject.NewProject(rancherProject)
	host1 = &rancherClient.Host{
		Name:  "Test1",
		State: "inactive",
	}
	project.AddHost(host1)
	host2 = &rancherClient.Host{
		Name:  "Test1",
		State: "disconnected",
	}
	project.AddHost(host2)
	monitoringData, err = CheckHostsProject(project)
	assert.Equal(t, 2, monitoringData.Status())
	assert.Equal(t, 2, len(monitoringData.Perfdatas()))

	// When project is disabled
	rancherProject = &rancherClient.Project{
		Name:  "Test",
		State: "inactive",
	}
	project, _ = modelProject.NewProject(rancherProject)
	monitoringData, err = CheckHostsProject(project)
	assert.Equal(t, 0, monitoringData.Status())
	assert.Equal(t, 0, len(monitoringData.Perfdatas()))
}
