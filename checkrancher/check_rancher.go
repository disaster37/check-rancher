package checkrancher

import (
	"time"

	"github.com/disaster37/go-nagios"
	"github.com/pkg/errors"
	rancher "github.com/rancher/go-rancher/v2"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// CheckRancher is implementation of MonitorRancher
type CheckRancher struct {
	client *rancher.RancherClient
}

// MonitorES is interface of elasticsearch monitoring
type MonitorRancher interface {
	CheckCertificates(projectName string, thresholdWarningDay int) (*nagiosPlugin.Monitoring, error)
	CheckHosts(projectName string) (*nagiosPlugin.Monitoring, error)
	CheckStack(projectName string, stackName string) (*nagiosPlugin.Monitoring, error)
}

func manageRancherGlobalParameters(c *cli.Context) (MonitorRancher, error) {

	if c.String("url") == "" {
		return nil, errors.New("You must set --url parameter")
	}

	if c.String("rancher-key") == "" {
		return nil, errors.New("You must set --rancher-key parameter")
	}

	if c.String("rancher-secret") == "" {
		return nil, errors.New("You must set --rancher-secret parameter")
	}

	if c.String("project-name") == "" {
		return nil, errors.New("You must set --project-name parameter")
	}

	return NewCheckRancher(c.String("url"), c.String("rancher-key"), c.String("rancher-secret"))

}

//NewCheckES permit to initialize connexion on Elasticsearch cluster
func NewCheckRancher(URL string, key string, secret string) (MonitorRancher, error) {

	if URL == "" {
		return nil, errors.New("URL can't be empty")
	}
	if key == "" {
		return nil, errors.New("Key can't be empty")
	}
	if secret == "" {
		return nil, errors.New("Secret can't be empty")
	}
	log.Debugf("URL: %s", URL)
	log.Debugf("Key: %s", key)
	log.Debugf("Secret: xxx")

	checkRancher := &CheckRancher{}
	opts := &rancher.ClientOpts{
		Url:       URL,
		AccessKey: key,
		SecretKey: secret,
		Timeout:   time.Second * 5,
	}
	client, err := rancher.NewRancherClient(opts)
	if err != nil {
		return nil, err
	}

	checkRancher.client = client
	return checkRancher, nil
}
