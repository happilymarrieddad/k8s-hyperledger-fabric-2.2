package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Read returns the value at id in the world state
func (rc *ResourceTypesContract) Read(ctx contractapi.TransactionContextInterface, id string) (ret *ResourceType, err error) {
	resultsIterator, _, err := ctx.GetStub().GetQueryResultWithPagination(`{"selector": {"id":"`+id+`"}}`, 0, "")
	if err != nil {
		return
	}
	defer resultsIterator.Close()

	if resultsIterator.HasNext() {
		ret = new(ResourceType)
		queryResponse, err2 := resultsIterator.Next()
		if err2 != nil {
			return nil, err2
		}

		if err = json.Unmarshal(queryResponse.Value, ret); err != nil {
			return
		}
	} else {
		return nil, fmt.Errorf("Unable to find item in world state")
	}

	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, err
	}

	valAsBytes, err := ctx.GetStub().GetPrivateData(
		mspID+"ResourceTypesPrivateData",
		id,
	)
	if err == nil && len(valAsBytes) > 0 {
		ret.PrivateName = string(valAsBytes)
	}

	return
}
