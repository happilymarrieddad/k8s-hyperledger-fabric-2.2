package resources

import (
	"encoding/json"

	"k8s-hyperledger-fabric-2.2/go-api/hyperledger"
	"k8s-hyperledger-fabric-2.2/go-api/models"
)

func Index(clients *hyperledger.Clients) (resources *models.Resources, err error) {
	resources = new(models.Resources)

	res, err := clients.Query("ibm", "mainchannel", "resources", "index", [][]byte{
		[]byte("{\"selector\":{ \"active\": { \"$eq\":true } }}"),
	})
	if err != nil {
		return
	}

	if err = json.Unmarshal(res, resources); err != nil {
		return
	}

	return
}
