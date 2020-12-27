package resources

import (
	"encoding/json"
	"os"

	"k8s-hyperledger-fabric-2.2/go-api/hyperledger"
	"k8s-hyperledger-fabric-2.2/go-api/models"
)

type UpdateOpts struct {
	Replace bool
}

func Update(clients *hyperledger.Clients, id string, rr *models.Resource, opts *UpdateOpts) (*models.Resource, error) {
	if !opts.Replace {
		existingResource, err := Show(clients, id)
		if err != nil {
			return nil, err
		}

		if rr.Name != "" && existingResource.Name != rr.Name {
			existingResource.Name = rr.Name
		}

		if rr.ResourceTypeID != "" && existingResource.ResourceTypeID != rr.ResourceTypeID {
			existingResource.ResourceTypeID = rr.ResourceTypeID
		}
	}

	rr.Active = true

	packet, err := json.Marshal(rr)
	if err != nil {
		return nil, err
	}

	MSPID := os.Getenv("HYPERLEDGER_MSP_ID")
	if len(MSPID) == 0 {
		MSPID = "ibm"
	}

	if _, err = clients.Invoke(MSPID, "mainchannel", "resources", "update", [][]byte{
		[]byte(id),
		packet,
	}); err != nil {
		return nil, err
	}

	return rr, nil
}
