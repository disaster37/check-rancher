package rancherapi

import (
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func (t *APITestSuite) TestClusterList() {

	// Normal use case
	clusterResp := &ClustersResponse{
		Data: []*Cluster{
			{
				Name: "test",
				ID:   "test_id",
			},
			{
				Name: "test2",
				ID:   "test2_id",
			},
		},
	}

	responder, err := httpmock.NewJsonResponder(200, clusterResp)
	if err != nil {
		panic(err)
	}
	httpmock.RegisterResponder("GET", "http://localhost/v3/clusters", responder)
	clusters, err := t.client.Cluster().List()
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), 2, len(clusters))

	// When not found
	httpmock.RegisterResponder("GET", "http://localhost/v3/clusters",
		httpmock.NewStringResponder(404, ""))
	_, err = t.client.Cluster().List()
	assert.Error(t.T(), err)

}

func (t *APITestSuite) TestClusterGetByName() {

	// Normal use case
	clusterResp := &ClustersResponse{
		Data: []*Cluster{
			{
				Name: "test",
				ID:   "test_id",
			},
		},
	}
	responder, err := httpmock.NewJsonResponder(200, clusterResp)
	if err != nil {
		panic(err)
	}
	httpmock.RegisterResponder("GET", "http://localhost/v3/clusters?name=test", responder)
	cluster, err := t.client.Cluster().GetByName("test")
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), "test_id", cluster.ID)

	// When empty list
	clusterResp = &ClustersResponse{
		Data: []*Cluster{},
	}
	responder, err = httpmock.NewJsonResponder(200, clusterResp)
	if err != nil {
		panic(err)
	}
	httpmock.RegisterResponder("GET", "http://localhost/v3/clusters?name=test", responder)
	cluster, err = t.client.Cluster().GetByName("test")
	assert.NoError(t.T(), err)
	assert.Nil(t.T(), cluster)

	// When not found from API
	httpmock.RegisterResponder("GET", "http://localhost/v3/clusters?name=test",
		httpmock.NewStringResponder(404, ""))
	_, err = t.client.Cluster().GetByName("test")
	assert.Error(t.T(), err)

}
