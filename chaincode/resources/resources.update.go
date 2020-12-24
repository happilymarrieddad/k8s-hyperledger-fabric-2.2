package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Update changes the value with id in the world state
func (rc *ResourcesContract) Update(ctx contractapi.TransactionContextInterface, id string, name string, resourceTypeID string, newOwner string) error {
	existing, err := ctx.GetStub().GetState(id)

	if err != nil {
		return fmt.Errorf("Unable to interact with world state")
	}

	if existing == nil {
		return fmt.Errorf("Cannot update world state pair with id %s. Does not exist", id)
	}

	var existingResource *Resource
	if err = json.Unmarshal(existing, &existingResource); err != nil {
		return fmt.Errorf("Unable to unmarshal existing into object")
	}
	if len(name) > 0 {
		existingResource.Name = name
	}
	if len(resourceTypeID) > 0 {
		existingResource.ResourceTypeID = resourceTypeID
	}

	// Change ownership if needed
	if len(newOwner) > 0 {
		existingResource.Owner = newOwner
	}

	newValue, err := json.Marshal(existingResource)
	if err != nil {
		return fmt.Errorf("Unable to marshal new object")
	}

	if err = ctx.GetStub().PutState(id, newValue); err != nil {
		return fmt.Errorf("Unable to interact with world state")
	}

	return nil
}
