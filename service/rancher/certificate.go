package rancherService

import (
	rancherClient "github.com/rancher/go-rancher/v2"
	log "github.com/sirupsen/logrus"
)

// Permit to retrieve all certificates
func (r *Client) GetCertificates() ([]rancherClient.Certificate, error) {
	log.Debugf("Looking up all Rancher certificate")

	certificates, err := r.client.Certificate.List(nil)

	if err != nil {
		return nil, err
	}

	if len(certificates.Data) == 0 {
		return nil, nil
	}

	log.Debugf("Found %d Rancher certificates", len(certificates.Data))
	return certificates.Data, nil
}
