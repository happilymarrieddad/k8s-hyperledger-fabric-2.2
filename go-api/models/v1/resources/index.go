package resources

import (
	"encoding/json"
	"os"

	"k8s-hyperledger-fabric-2.2/go-api/hyperledger"
	"k8s-hyperledger-fabric-2.2/go-api/models"
)

func Index(clients *hyperledger.Clients) (resources *models.Resources, err error) {
	resources = new(models.Resources)

	MSPID := os.Getenv("HYPERLEDGER_MSP_ID")
	if len(MSPID) == 0 {
		MSPID = "ibm"
	}

	res, err := clients.Query(MSPID, "mainchannel", "resources", "index", [][]byte{
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
