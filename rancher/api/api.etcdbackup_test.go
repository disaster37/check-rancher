package rancherapi

import (
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func (t *APITestSuite) TestETCDBackupList() {

	// Normal use case
	backupResp := &ETCDBackupResp{
		Data: []*ETCDBackup{
			{
				ID: "test1",
			},
			{
				ID: "test2",
			},
		},
	}

	responder, err := httpmock.NewJsonResponder(200, backupResp)
	if err != nil {
		panic(err)
	}
	httpmock.RegisterResponder("GET", "http://localhost/v3/etcdbackups", responder)
	clusters, err := t.client.ETCDBackup().List()
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), 2, len(clusters))

	// When not found
	httpmock.RegisterResponder("GET", "http://localhost/v3/etcdbackups",
		httpmock.NewStringResponder(404, ""))
	_, err = t.client.Cluster().List()
	assert.Error(t.T(), err)

}

func (t *APITestSuite) TestETCDBackupByClusterID() {

	// Normal use case
	backupResp := &ETCDBackupResp{
		Data: []*ETCDBackup{
			{
				ID: "test1",
			},
			{
				ID: "test2",
			},
		},
	}
	responder, err := httpmock.NewJsonResponder(200, backupResp)
	if err != nil {
		panic(err)
	}
	httpmock.RegisterResponder("GET", "http://localhost/v3/etcdbackups?clusterId=test", responder)
	backups, err := t.client.ETCDBackup().ListByClusterID("test")
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), 2, len(backups))

	// When empty list
	backupResp = &ETCDBackupResp{
		Data: []*ETCDBackup{},
	}
	responder, err = httpmock.NewJsonResponder(200, backupResp)
	if err != nil {
		panic(err)
	}
	httpmock.RegisterResponder("GET", "http://localhost/v3/etcdbackups?clusterId=test", responder)
	backups, err = t.client.ETCDBackup().ListByClusterID("test")
	assert.NoError(t.T(), err)
	assert.Empty(t.T(), backups)

	// When not found from API
	httpmock.RegisterResponder("GET", "http://localhost/v3/etcdbackups?clusterId=test",
		httpmock.NewStringResponder(404, ""))
	_, err = t.client.ETCDBackup().ListByClusterID("test")
	assert.Error(t.T(), err)

}
