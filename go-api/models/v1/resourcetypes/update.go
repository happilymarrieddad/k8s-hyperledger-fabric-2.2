package resourcetypes

import (
	"errors"

	"k8s-hyperledger-fabric-2.2/go-api/models"
)

type UpdateOpts struct {
	Replace bool
}

func Update(id string, usr *models.ResourceType, opts *UpdateOpts) (*models.ResourceType, error) {
	var exists bool

	if opts == nil {
		opts = &UpdateOpts{}
	}

	for index, resourcetype := range mockResourceTypes {
		if resourcetype.ID == id {

			if opts.Replace {
				usr.ID = mockResourceTypes[index].ID
				mockResourceTypes[index] = *usr
			} else {
				if len(usr.Name) > 0 {
					mockResourceTypes[index].Name = usr.Name
				}
			}
			exists = true
			usr = &mockResourceTypes[index]
		}
	}

	if !exists {
		return nil, errors.New("unable to update resourcetype because resourcetype was not found")
	}

	return usr, nil
}
