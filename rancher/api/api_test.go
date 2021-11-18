package rancherapi

import (
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"
)

type APITestSuite struct {
	suite.Suite
	client API
}

func TestAPISuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}

func (t *APITestSuite) SetupTest() {
	restyClient := resty.New().
		SetHostURL("http://localhost/v3").
		SetHeader("Content-Type", "application/json")
	httpmock.ActivateNonDefault(restyClient.GetClient())

	t.client = New(restyClient)
}

func (t *APITestSuite) BeforeTest(suiteName, testName string) {
	httpmock.Reset()
}
