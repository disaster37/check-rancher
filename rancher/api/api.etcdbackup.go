package rancherapi

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	basePathETCDBackup = "etcdbackups"
)

type ETCDBackupAPIImpl struct {
	client *resty.Client
}

func NewETCDBackupAPI(client *resty.Client) ETCDBackupAPI {
	return &ETCDBackupAPIImpl{
		client: client,
	}
}

type ETCDBackup struct {
	ID            string    `json:"id,omitempty"`
	ClusterID     string    `json:"clusterId,omitempty"`
	CreatedAt     time.Time `json:"created,omitempty"`
	Filename      string    `json:"filename,omitempty"`
	State         string    `json:"state,omitempty"`
	Transitioning string    `json:"transitioning,omitempty"`
}

type ETCDBackupResp struct {
	Data []*ETCDBackup `json:"data,omitempty"`
}

func (api *ETCDBackupAPIImpl) List() ([]*ETCDBackup, error) {

	path := fmt.Sprintf(basePathETCDBackup)

	resp, err := api.client.R().
		Get(path)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() >= 300 {
		return nil, errors.Errorf("Error when list ETCD backups: %s", resp.Body())
	}

	respBackups := new(ETCDBackupResp)
	if err := json.Unmarshal(resp.Body(), respBackups); err != nil {
		return nil, err
	}

	log.Debugf("Backups: %+v", respBackups.Data)

	return respBackups.Data, nil
}

func (api *ETCDBackupAPIImpl) ListByClusterID(clusterID string) ([]*ETCDBackup, error) {

	path := fmt.Sprintf(basePathETCDBackup)

	resp, err := api.client.R().
		SetQueryParam("clusterId", clusterID).
		Get(path)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() >= 300 {
		return nil, errors.Errorf("Error when list ETCD backups for cluster %s: %s", clusterID, resp.Body())
	}

	respBackup := new(ETCDBackupResp)
	if err := json.Unmarshal(resp.Body(), respBackup); err != nil {
		return nil, err
	}

	log.Debugf("Backups for cluster %s: %+v", clusterID, respBackup.Data)

	return respBackup.Data, nil
}
