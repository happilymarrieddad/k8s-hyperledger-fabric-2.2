package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Update changes the value with id in the world state
func (rc *ResourceTypesContract) Update(ctx contractapi.TransactionContextInterface, id string, name string) error {
	existing, err := ctx.GetStub().GetState(id)

	if err != nil {
		return fmt.Errorf("Unable to interact with world state")
	}

	if existing == nil {
		return fmt.Errorf("Cannot update world state pair with id %s. Does not exist", id)
	}

	var existingResourceType *ResourceType
	if err = json.Unmarshal(existing, &existingResourceType); err != nil {
		return fmt.Errorf("Unable to unmarshal existing into object")
	}
	existingResourceType.Name = name

	newValue, err := json.Marshal(existingResourceType)
	if err != nil {
		return fmt.Errorf("Unable to marshal new object")
	}

	if err = ctx.GetStub().PutState(id, newValue); err != nil {
		return fmt.Errorf("Unable to interact with world state")
	}

	return nil
}
