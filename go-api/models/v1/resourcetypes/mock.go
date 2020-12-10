package resourcetypes

import (
	"k8s-hyperledger-fabric-2.2/go-api/models"
)

var mockResourceTypes models.ResourceTypes

func init() {
	iron, _ := models.NewResourceType("Iron")
	copper, _ := models.NewResourceType("Copper")
	platinum, _ := models.NewResourceType("Platinum")

	mockResourceTypes = models.ResourceTypes{
		*iron,
		*copper,
		*platinum,
	}
}
