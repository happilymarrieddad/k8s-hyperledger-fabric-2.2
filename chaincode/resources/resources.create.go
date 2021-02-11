package main

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric/common/util"
)

// Create adds a new id with value to the world state
func (rc *ResourcesContract) Create(ctx contractapi.TransactionContextInterface, id, name, resourceTypeID string) error {

	// TODO: Verify this name is unique
	newResource := &Resource{
		ID:             id,
		Name:           name,
		ResourceTypeID: resourceTypeID,
		Active:         true,
	}

	if len(newResource.ID) == 0 {
		newResource.ID = uuid.New().String()
	}

	// This shouldn't ever happen but just in case
	existing, err := ctx.GetStub().GetState(newResource.ID)
	if err != nil {
		return fmt.Errorf("Unable to interact with world state")
	}

	if existing != nil {
		return fmt.Errorf("Cannot create world state pair with id %s. Already exists", newResource.ID)
	}

	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("Unable to interact with client identity")
	}
	newResource.Owner = mspID

	chainCodeArgs := util.ToChaincodeArgs("Read", newResource.ResourceTypeID)

	if res := ctx.GetStub().InvokeChaincode("resource_types", chainCodeArgs, ""); res.Status != 200 {
		return fmt.Errorf("Resource type '%s' does not exist", newResource.ResourceTypeID)
	}

	bytes, err := json.Marshal(newResource)
	if err != nil {
		return fmt.Errorf("Unable to marshal object")
	}

	if err = ctx.GetStub().PutState(newResource.ID, bytes); err != nil {
		return fmt.Errorf("Unable to interact with world state")
	}

	return nil
}
