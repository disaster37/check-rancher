package cmd

import (
	"time"

	"github.com/disaster37/check-rancher/v2/rancher"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func getClientWrapper(c *cli.Context) (*rancher.Client, error) {

	url := c.String("url")
	accessKey := c.String("access-key")
	secretKey := c.String("secret-key")
	disableVerifySSL := c.Bool("self-signed-certificate")
	listCAPath := c.StringSlice("ca-path")
	timeout := c.Duration("timeout")
	debug := c.Bool("debug")

	log.Debugf("URL: %s", url)
	log.Debugf("Access key: %s", accessKey)
	log.Debugf("Self signed certificate: %t", disableVerifySSL)
	log.Debugf("List CA path: %+v", listCAPath)
	log.Debugf("Debug: %t", debug)

	return getClient(url, accessKey, secretKey, disableVerifySSL, listCAPath, timeout, debug)

}

func getClient(url string, accessKey string, secretKey string, disableVerifySSL bool, listCAPath []string, timeout time.Duration, debug bool) (*rancher.Client, error) {

	if url == "" {
		return nil, errors.New("You need to set url")
	}

	config := rancher.Config{
		Address:          url,
		AccessKey:        accessKey,
		SecretKey:        secretKey,
		DisableVerifySSL: disableVerifySSL,
		CAs:              listCAPath,
		Timeout:          timeout,
		Debug:            debug,
	}

	return rancher.NewClient(config)
}
