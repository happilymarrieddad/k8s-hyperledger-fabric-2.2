package resources

import (
	"encoding/json"
	"os"

	"k8s-hyperledger-fabric-2.2/go-api/hyperledger"
	"k8s-hyperledger-fabric-2.2/go-api/models"
)

func Store(clients *hyperledger.Clients, name string, typeID string) (resource *models.Resource, err error) {
	resource = new(models.Resource)

	resource.Name = name
	resource.ResourceTypeID = typeID
	resource.Active = true

	MSPID := os.Getenv("HYPERLEDGER_MSP_ID")
	if len(MSPID) == 0 {
		MSPID = "ibm"
	}

	res, err := clients.Invoke(MSPID, "mainchannel", "resources", "create", [][]byte{
		[]byte(""),
		[]byte(name),
		[]byte(typeID),
	})
	if err != nil {
		return
	}

	if err = json.Unmarshal(res, resource); err != nil {
		return
	}

	return
}
