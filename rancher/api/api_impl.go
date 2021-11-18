package rancherapi

import "github.com/go-resty/resty/v2"

type APIImpl struct {
	cluster    ClusterAPI
	etcdBackup ETCDBackupAPI
}

func New(client *resty.Client) API {
	return &APIImpl{
		cluster:    NewClusterAPI(client),
		etcdBackup: NewETCDBackupAPI(client),
	}
}

func (api *APIImpl) Cluster() ClusterAPI {
	return api.cluster
}

func (api *APIImpl) ETCDBackup() ETCDBackupAPI {
	return api.etcdBackup
}
