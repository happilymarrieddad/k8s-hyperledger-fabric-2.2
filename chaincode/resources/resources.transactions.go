package main

import (
	"encoding/json"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Transactions get all the transactions of an id
func (rc *ResourcesContract) Transactions(
	ctx contractapi.TransactionContextInterface,
	id string,
) ([]*ResourceTransactionItem, error) {
	historyIface, err := ctx.GetStub().GetHistoryForKey(id)
	if err != nil {
		return nil, err
	}

	var rets []*ResourceTransactionItem
	for historyIface.HasNext() {
		val, err := historyIface.Next()
		if err != nil {
			return nil, err
		}

		var res Resource
		if err = json.Unmarshal(val.Value, &res); err != nil {
			return nil, err
		}

		rets = append(rets, &ResourceTransactionItem{
			TXID:      val.TxId,
			Timestamp: int64(val.Timestamp.GetNanos()),
			Resource:  res,
		})
	}

	return rets, nil
}
