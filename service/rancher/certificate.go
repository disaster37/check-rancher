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

// Permit to retrive hosts linked to environment
func (r *Client) FindCertificatesByProjectLink(project *rancherClient.Project) ([]rancherClient.Certificate, error) {
	log.Debugf("Looking up Rancher certificates for project: %s", project.Name)
	certificateCollection := rancherClient.CertificateCollection{}
	err := r.client.GetLink(project.Resource, "certificates", &certificateCollection)
	if err != nil {
		return nil, err
	}

	if len(certificateCollection.Data) == 0 {
		return nil, nil
	}

	log.Debugf("Found %d certificates link with project %s", len(certificateCollection.Data), project.Name)

	return certificateCollection.Data, nil
}
