package main

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric/common/util"
)

// Create adds a new id with value to the world state
func (rc *ResourcesContract) Create(ctx contractapi.TransactionContextInterface, name string, resourceTypeID string) error {
	id := uuid.New().String()

	// This shouldn't ever happen but just in case
	existing, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("Unable to interact with world state")
	}

	if existing != nil {
		return fmt.Errorf("Cannot create world state pair with id %s. Already exists", id)
	}

	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("Unable to interact with client identity")
	}

	// TODO: Verify this name is unique
	newResource := &Resource{
		ID:             id,
		Name:           name,
		ResourceTypeID: resourceTypeID,
		Active:         true,
		Owner:          mspID,
	}

	chainCodeArgs := util.ToChaincodeArgs("Read", resourceTypeID)

	if res := ctx.GetStub().InvokeChaincode("resource_types", chainCodeArgs, ""); res.Status != 200 {
		return fmt.Errorf("Resource type '%s' does not exist", resourceTypeID)
	}

	bytes, err := json.Marshal(newResource)
	if err != nil {
		return fmt.Errorf("Unable to marshal object")
	}

	if err = ctx.GetStub().PutState(id, bytes); err != nil {
		return fmt.Errorf("Unable to interact with world state")
	}

	return nil
}
