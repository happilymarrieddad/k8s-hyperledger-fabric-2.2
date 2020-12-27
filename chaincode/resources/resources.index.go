package main

import (
	"encoding/json"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Index - read all resources from the world state
func (rc *ResourcesContract) Index(
	ctx contractapi.TransactionContextInterface,
) (rets []*Resource, err error) {
	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return
	}

	resultsIterator, _, err := ctx.GetStub().GetQueryResultWithPagination(`{"selector": {"id":{"$ne":"-"},"owner":"`+mspID+`"}}`, 0, "")
	if err != nil {
		return
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		queryResponse, err2 := resultsIterator.Next()
		if err2 != nil {
			return nil, err2
		}

		res := new(Resource)
		if err = json.Unmarshal(queryResponse.Value, res); err != nil {
			return
		}

		rets = append(rets, res)
	}

	return
}
