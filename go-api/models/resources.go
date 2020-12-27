package models

import "time"

type ResourceTypes []ResourceType

type ResourceType struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Active      bool   `json:"active"`
	PrivateName string `json:"privateName"`
}

func NewResourceType(name string) (resourceType *ResourceType, err error) {
	resourceType = new(ResourceType)

	if resourceType.ID, err = genUUID(); err != nil {
		return
	}

	resourceType.Name = name
	resourceType.Active = true

	return
}

type Resources []Resource

type Resource struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	ResourceTypeID string `json:"resource_type_id"`
	Active         bool   `json:"active"`
	Owner          string `json:"owner"`
}

func NewResource(name string, typeId string, weight int, arrivalTime *time.Time) (resource *Resource, err error) {
	resource = new(Resource)

	if resource.ID, err = genUUID(); err != nil {
		return
	}

	resource.Name = name
	resource.ResourceTypeID = typeId
	resource.Active = true

	return
}
