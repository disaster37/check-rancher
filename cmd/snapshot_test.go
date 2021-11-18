package cmd

import (
	"errors"
	"time"

	rancherapi "github.com/disaster37/check-rancher/v2/rancher/api"
	nagiosPlugin "github.com/disaster37/go-nagios"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func (t CmdTestSuite) TestCheckSnapshot() {

	t.mockClient.EXPECT().Cluster().AnyTimes().Return(t.mockCluster)
	t.mockClient.EXPECT().ETCDBackup().AnyTimes().Return(t.mockETCDBackup)

	// Normale use case when cluster found and snapshot found and Ok
	cluster := &rancherapi.Cluster{
		Name: "test",
		ID:   "test",
	}
	t.mockCluster.
		EXPECT().
		GetByName(gomock.Eq("test")).
		Return(cluster, nil)
	backups := []*rancherapi.ETCDBackup{
		{
			ID:            "test",
			ClusterID:     "test",
			CreatedAt:     time.Now(),
			State:         "active",
			Transitioning: "no",
		},
	}
	t.mockETCDBackup.
		EXPECT().
		ListByClusterID(gomock.Eq("test")).
		Return(backups, nil)
	md := checkSnapshot("test", 24*time.Hour, t.client)
	assert.Equal(t.T(), nagiosPlugin.STATUS_OK, md.Status())

	// Normale use case when cluster found and snapshot found and Ko
	cluster = &rancherapi.Cluster{
		Name: "test",
		ID:   "test",
	}
	t.mockCluster.
		EXPECT().
		GetByName(gomock.Eq("test")).
		Return(cluster, nil)
	backups = []*rancherapi.ETCDBackup{
		{
			ID:            "test",
			ClusterID:     "test",
			CreatedAt:     time.Now(),
			State:         "failed",
			Transitioning: "no",
		},
	}
	t.mockETCDBackup.
		EXPECT().
		ListByClusterID(gomock.Eq("test")).
		Return(backups, nil)
	md = checkSnapshot("test", 24*time.Hour, t.client)
	assert.Equal(t.T(), nagiosPlugin.STATUS_CRITICAL, md.Status())

	// Normale use case when cluster found and multiples snapshots. It take the last one
	cluster = &rancherapi.Cluster{
		Name: "test",
		ID:   "test",
	}
	t.mockCluster.
		EXPECT().
		GetByName(gomock.Eq("test")).
		Return(cluster, nil)
	backups = []*rancherapi.ETCDBackup{
		{
			ID:            "test",
			ClusterID:     "test",
			CreatedAt:     time.Now(),
			State:         "failed",
			Transitioning: "no",
		},
		{
			ID:            "test",
			ClusterID:     "test",
			CreatedAt:     time.Now().Add(-1 * time.Hour),
			State:         "ok",
			Transitioning: "no",
		},
	}
	t.mockETCDBackup.
		EXPECT().
		ListByClusterID(gomock.Eq("test")).
		Return(backups, nil)
	md = checkSnapshot("test", 24*time.Hour, t.client)
	assert.Equal(t.T(), nagiosPlugin.STATUS_CRITICAL, md.Status())

	// Normale use case when cluster found and snapshot found and OK but older than
	cluster = &rancherapi.Cluster{
		Name: "test",
		ID:   "test",
	}
	t.mockCluster.
		EXPECT().
		GetByName(gomock.Eq("test")).
		Return(cluster, nil)
	backups = []*rancherapi.ETCDBackup{
		{
			ID:            "test",
			ClusterID:     "test",
			CreatedAt:     time.Now().Add(-72 * time.Hour),
			State:         "active",
			Transitioning: "no",
		},
	}
	t.mockETCDBackup.
		EXPECT().
		ListByClusterID(gomock.Eq("test")).
		Return(backups, nil)
	md = checkSnapshot("test", 24*time.Hour, t.client)
	assert.Equal(t.T(), nagiosPlugin.STATUS_CRITICAL, md.Status())

	// Normale use case when cluster found and snapshot found and OK on transitioning
	cluster = &rancherapi.Cluster{
		Name: "test",
		ID:   "test",
	}
	t.mockCluster.
		EXPECT().
		GetByName(gomock.Eq("test")).
		Return(cluster, nil)
	backups = []*rancherapi.ETCDBackup{
		{
			ID:            "test",
			ClusterID:     "test",
			CreatedAt:     time.Now(),
			State:         "active",
			Transitioning: "yes",
		},
	}
	t.mockETCDBackup.
		EXPECT().
		ListByClusterID(gomock.Eq("test")).
		Return(backups, nil)
	md = checkSnapshot("test", 24*time.Hour, t.client)
	assert.Equal(t.T(), nagiosPlugin.STATUS_CRITICAL, md.Status())

	// Normale use case when cluster found and no snapshot
	cluster = &rancherapi.Cluster{
		Name: "test",
		ID:   "test",
	}
	t.mockCluster.
		EXPECT().
		GetByName(gomock.Eq("test")).
		Return(cluster, nil)
	backups = []*rancherapi.ETCDBackup{}
	t.mockETCDBackup.
		EXPECT().
		ListByClusterID(gomock.Eq("test")).
		Return(backups, nil)
	md = checkSnapshot("test", 24*time.Hour, t.client)
	assert.Equal(t.T(), nagiosPlugin.STATUS_CRITICAL, md.Status())

	// Use case when cluster found and snapshot error
	cluster = &rancherapi.Cluster{
		Name: "test",
		ID:   "test",
	}
	t.mockCluster.
		EXPECT().
		GetByName(gomock.Eq("test")).
		Return(cluster, nil)
	backups = []*rancherapi.ETCDBackup{}
	t.mockETCDBackup.
		EXPECT().
		ListByClusterID(gomock.Eq("test")).
		Return(backups, errors.New("fake error"))
	md = checkSnapshot("test", 24*time.Hour, t.client)
	assert.Equal(t.T(), nagiosPlugin.STATUS_UNKNOWN, md.Status())

	// Usecase when cluster not found
	cluster = nil
	t.mockCluster.
		EXPECT().
		GetByName(gomock.Eq("test")).
		Return(cluster, nil)
	md = checkSnapshot("test", 24*time.Hour, t.client)
	assert.Equal(t.T(), nagiosPlugin.STATUS_UNKNOWN, md.Status())

	// Usecase when cluster error
	cluster = nil
	t.mockCluster.
		EXPECT().
		GetByName(gomock.Eq("test")).
		Return(cluster, errors.New("fake error"))
	md = checkSnapshot("test", 24*time.Hour, t.client)
	assert.Equal(t.T(), nagiosPlugin.STATUS_UNKNOWN, md.Status())

}
