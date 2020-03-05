package checkrancher

import (
	"github.com/pkg/errors"
	rancher "github.com/rancher/go-rancher/v2"
	log "github.com/sirupsen/logrus"
)

// Permit to retrieve an existing project by name
func (h *CheckRancher) FindProjectByName(name string) (*rancher.Project, error) {

	if name == "" {
		return nil, errors.New("Name can't be empty")
	}
	log.Debugf("Looking up Rancher project by name: %s", name)

	projects, err := h.client.Project.List(&rancher.ListOpts{
		Filters: map[string]interface{}{
			"name":         name,
			"removed_null": nil,
		},
	})

	if err != nil {
		return nil, err
	}

	if len(projects.Data) == 0 {
		return nil, nil
	}

	log.Debugf("Found existing Rancher project by name: %s", name)
	return &projects.Data[0], nil
}
