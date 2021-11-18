package rancher

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/disaster37/check-rancher/v2/rancher/api"
	"github.com/go-resty/resty/v2"
)

// Config contain the value to access on Kibana API
type Config struct {
	Address          string
	AccessKey        string
	SecretKey        string
	DisableVerifySSL bool
	CAs              []string
	Timeout          time.Duration
	Debug            bool
}

// Client contain the REST client and the API specification
type Client struct {
	API rancherapi.API
}

// NewDefaultClient init client with empty config
func NewDefaultClient() (*Client, error) {
	return NewClient(Config{})
}

// NewClient init client with custom config
func NewClient(cfg Config) (*Client, error) {
	if cfg.Address == "" {
		cfg.Address = "https://localhost/v3"
	}

	restyClient := resty.New().
		SetHostURL(fmt.Sprintf("%s/v3", cfg.Address)).
		SetBasicAuth(cfg.AccessKey, cfg.SecretKey).
		SetHeader("Content-Type", "application/json").
		SetTimeout(cfg.Timeout).
		SetDebug(cfg.Debug).
		SetCookieJar(nil)

	for _, path := range cfg.CAs {
		restyClient.SetRootCertificate(path)
	}

	if cfg.DisableVerifySSL == true {
		restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	client := &Client{
		API: rancherapi.New(restyClient),
	}

	return client, nil

}
