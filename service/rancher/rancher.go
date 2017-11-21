package rancherService

import (
	rancherClient "github.com/rancher/go-rancher/v2"
	"time"
    log "github.com/sirupsen/logrus"
)

type Client struct {
	client *rancherClient.RancherClient
}

// Permit to get Rancher API connection
func NewClient(rancherUrl string, rancherAccessKey string, rancherSecretKey string) (*Client, error) {
	opts := &rancherClient.ClientOpts{
		Url:       rancherUrl,
		AccessKey: rancherAccessKey,
		SecretKey: rancherSecretKey,
		Timeout:   time.Second * 5,
	}

	var err error
	var apiClient *rancherClient.RancherClient
	maxTime := 10 * time.Second

	for i := 1 * time.Second; i < maxTime; i *= time.Duration(2) {
		apiClient, err = rancherClient.NewRancherClient(opts)
		if err == nil {
			break
		}
		time.Sleep(i)
	}

	if err != nil {
		return nil, err
	}
	
	log.Debugf("Successfully connected to Rancher url: %s", rancherUrl)

	return &Client{apiClient}, nil
}
