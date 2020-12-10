package resourcetypes

import (
	"errors"
)

func Destroy(id string) error {
	var exists bool

	for index, resourcetype := range mockResourceTypes {
		if resourcetype.ID == id {
			mockResourceTypes = append(mockResourceTypes[:index], mockResourceTypes[index+1:]...)
			exists = true
		}
	}

	if !exists {
		return errors.New("unable to delete resourcetype because resourcetype was not found")
	}

	return nil
}
