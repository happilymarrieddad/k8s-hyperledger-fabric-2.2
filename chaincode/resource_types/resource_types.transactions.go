package main

import (
	"encoding/json"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Transactions get all the transactions of an id
func (rc *ResourceTypesContract) Transactions(
	ctx contractapi.TransactionContextInterface,
	id string,
) ([]*ResourceTypeTransactionItem, error) {
	historyIface, err := ctx.GetStub().GetHistoryForKey(id)
	if err != nil {
		return nil, err
	}

	var rets []*ResourceTypeTransactionItem
	for historyIface.HasNext() {
		val, err := historyIface.Next()
		if err != nil {
			return nil, err
		}

		var res ResourceTypeIndex
		if err = json.Unmarshal(val.Value, &res); err != nil {
			return nil, err
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
			res.PrivateData = ResourceTypePrivateData{
				ResourceTypeID: id,
				DisplayName:    string(valAsBytes),
				Active:         res.Active,
			}
		}

		rets = append(rets, &ResourceTypeTransactionItem{
			TXID:         val.TxId,
			Timestamp:    int64(val.Timestamp.GetNanos()),
			ResourceType: res,
		})
	}

	return rets, nil
}
