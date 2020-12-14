package resources

import (
	"encoding/json"
	"errors"

	"k8s-hyperledger-fabric-2.2/go-api/hyperledger"
	"k8s-hyperledger-fabric-2.2/go-api/models"
)

func Show(clients *hyperledger.Clients, id string) (resource *models.Resource, err error) {

	resources := new(models.Resources)

	res, err := clients.Query("ibm", "mainchannel", "resources", "queryString", [][]byte{
		[]byte("{\"selector\":{ \"id\": { \"$eq\":\"" + id + "\" } }}"),
	})
	if err != nil {
		return
	}

	if err = json.Unmarshal(res, resources); err != nil {
		return
	}

	list := *resources

	if len(list) == 0 {
		err = errors.New("unable to find resource with id " + id)
		return
	}

	resource = &list[0]

	return
}
