package resources

import (
	"encoding/json"

	uuid "github.com/satori/go.uuid"

	"k8s-hyperledger-fabric-2.2/go-api/hyperledger"
	"k8s-hyperledger-fabric-2.2/go-api/models"
)

func Store(clients *hyperledger.Clients, name string, typeID string) (resource *models.Resource, err error) {
	resource = new(models.Resource)

	resource.ID = uuid.NewV4().String()
	resource.Name = name
	resource.ResourceTypeID = typeID
	resource.Active = true

	packet, err := json.Marshal(resource)
	if err != nil {
		return
	}

	if _, err = clients.Invoke("ibm", "mainchannel", "resources", "store", [][]byte{
		packet,
	}); err != nil {
		return
	}

	return
}
