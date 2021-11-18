package cmd

import (
	"testing"
	"time"

	"github.com/disaster37/check-rancher/v2/rancher"
	"github.com/disaster37/check-rancher/v2/rancher/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CmdTestSuite struct {
	suite.Suite
	client         *rancher.Client
	mockClient     *mocks.MockAPI
	mockCluster    *mocks.MockClusterAPI
	mockETCDBackup *mocks.MockETCDBackupAPI
	mockCtrl       *gomock.Controller
}

func TestAPISuite(t *testing.T) {
	suite.Run(t, new(CmdTestSuite))
}

func (t *CmdTestSuite) BeforeTest(suiteName, testName string) {

	t.mockCtrl = gomock.NewController(t.T())
	t.mockClient = mocks.NewMockAPI(t.mockCtrl)
	t.mockCluster = mocks.NewMockClusterAPI(t.mockCtrl)
	t.mockETCDBackup = mocks.NewMockETCDBackupAPI(t.mockCtrl)

	t.client = &rancher.Client{
		API: t.mockClient,
	}

}

func (t *CmdTestSuite) AfterTest(suiteName, testName string) {
	defer t.mockCtrl.Finish()
}

func (t *CmdTestSuite) TestGetClient() {

	client, err := getClient("http://localhost", "user", "password", false, nil, 1*time.Second, true)
	assert.NoError(t.T(), err)
	assert.NotNil(t.T(), client)
}
