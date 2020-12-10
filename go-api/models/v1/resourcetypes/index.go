package resourcetypes

import "k8s-hyperledger-fabric-2.2/go-api/models"

func Index() (resourcetypes *models.ResourceTypes, err error) {
	resourcetypes = &mockResourceTypes

	return
}
