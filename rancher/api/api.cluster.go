package rancherapi

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	basePathCluster = "clusters"
)

type ClusterAPIImpl struct {
	client *resty.Client
}

func NewClusterAPI(client *resty.Client) ClusterAPI {
	return &ClusterAPIImpl{
		client: client,
	}
}

type Cluster struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type ClustersResponse struct {
	Data []*Cluster `json:"data,omitempty"`
}

func (api *ClusterAPIImpl) List() ([]*Cluster, error) {

	path := fmt.Sprintf(basePathCluster)

	resp, err := api.client.R().
		Get(path)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() >= 300 {
		return nil, errors.Errorf("Error when list clusters: %s", resp.Body())
	}

	respClusters := new(ClustersResponse)
	if err := json.Unmarshal(resp.Body(), respClusters); err != nil {
		return nil, err
	}

	log.Debugf("Clusters: %+v", respClusters.Data)

	return respClusters.Data, nil
}

func (api *ClusterAPIImpl) GetByName(name string) (*Cluster, error) {

	path := fmt.Sprintf(basePathCluster)

	resp, err := api.client.R().
		SetQueryParam("name", name).
		Get(path)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() >= 300 {
		return nil, errors.Errorf("Error when list clusters: %s", resp.Body())
	}

	respClusters := new(ClustersResponse)
	if err := json.Unmarshal(resp.Body(), respClusters); err != nil {
		return nil, err
	}

	if len(respClusters.Data) == 0 {
		return nil, nil
	}

	return respClusters.Data[0], nil
}
