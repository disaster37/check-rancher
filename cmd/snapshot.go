package cmd

import (
	"time"

	"github.com/disaster37/check-rancher/v2/rancher"
	rancherapi "github.com/disaster37/check-rancher/v2/rancher/api"
	"github.com/disaster37/go-nagios"
	"github.com/urfave/cli/v2"
)

func CheckSnapshot(c *cli.Context) error {
	md := nagiosPlugin.NewMonitoring()
	client, err := getClientWrapper(c)
	if err != nil {
		md.SetStatus(nagiosPlugin.STATUS_UNKNOWN)
		md.AddMessage(err.Error())
		md.ToSdtOut()
	}

	md = checkSnapshot(c.String("cluster-name"), c.Duration("max-older-than"), client)
	md.ToSdtOut()

	return nil

}

func checkSnapshot(clusterName string, maxDuration time.Duration, client *rancher.Client) *nagiosPlugin.Monitoring {
	md := nagiosPlugin.NewMonitoring()

	// Get cluster ID
	cluster, err := client.API.Cluster().GetByName(clusterName)
	if err != nil {
		md.SetStatus(nagiosPlugin.STATUS_UNKNOWN)
		md.AddMessage(err.Error())
		return md
	}
	if cluster == nil {
		md.SetStatus(nagiosPlugin.STATUS_UNKNOWN)
		md.AddMessage("Cluster %s not found", clusterName)
		return md
	}

	// Get cluster snapshots and select the last
	snapshots, err := client.API.ETCDBackup().ListByClusterID(cluster.ID)
	if err != nil {
		md.SetStatus(nagiosPlugin.STATUS_UNKNOWN)
		md.AddMessage(err.Error())
		return md
	}
	if len(snapshots) == 0 {
		md.SetStatus(nagiosPlugin.STATUS_CRITICAL)
		md.AddMessage("No ETCD backup for cluster %s", clusterName)
		return md
	}
	var lastSnapshot *rancherapi.ETCDBackup
	nbBackupFailed := 0
	for _, snapshot := range snapshots {
		if snapshot.Transitioning == "no" {
			if lastSnapshot == nil || snapshot.CreatedAt.After(lastSnapshot.CreatedAt) {
				lastSnapshot = snapshot
			}
			if snapshot.State != "active" {
				nbBackupFailed++
			}
		}
	}

	// Check state
	if lastSnapshot == nil {
		md.SetStatus(nagiosPlugin.STATUS_CRITICAL)
		md.AddMessage("No ETCD backup for cluster %s", clusterName)
	} else if lastSnapshot.State != "active" {
		md.SetStatus(nagiosPlugin.STATUS_CRITICAL)
		md.AddMessage("The lastETCD backup for cluster %s is not on active state (%s)", clusterName, lastSnapshot.State)
	} else if time.Now().Sub(lastSnapshot.CreatedAt) > maxDuration {
		md.SetStatus(nagiosPlugin.STATUS_CRITICAL)
		md.AddMessage("The lastETCD backup for cluster %s is older than %d hours (%s)", clusterName, int64(maxDuration.Hours()), lastSnapshot.CreatedAt)
	} else {
		md.SetStatus(nagiosPlugin.STATUS_OK)
		md.AddMessage("The lastETCD backup for cluster %s is ok (%s)", clusterName, lastSnapshot.CreatedAt)
	}

	md.AddPerfdata("nbBackupFailed", nbBackupFailed, "")
	md.AddPerfdata("nbBackupOk", len(snapshots)-nbBackupFailed, "")

	return md
}
