package main

import (
	"encoding/json"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Index - read all resources from the world state
func (rc *ResourceTypesContract) Index(
	ctx contractapi.TransactionContextInterface,
) (rets []*ResourceType, err error) {
	resultsIterator, _, err := ctx.GetStub().GetQueryResultWithPagination(`{"selector": {"id":{"$ne":"-"}}}`, 0, "")
	if err != nil {
		return
	}
	defer resultsIterator.Close()

	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return
	}

	for resultsIterator.HasNext() {
		queryResponse, err2 := resultsIterator.Next()
		if err2 != nil {
			return nil, err2
		}

		res := new(ResourceType)
		if err = json.Unmarshal(queryResponse.Value, res); err != nil {
			return
		}

		valAsBytes, err := ctx.GetStub().GetPrivateData(
			mspID+"ResourceTypesPrivateData",
			res.ID,
		)
		if err == nil && len(valAsBytes) > 0 {
			res.PrivateName = string(valAsBytes)
		}

		rets = append(rets, res)
	}

	return
}
