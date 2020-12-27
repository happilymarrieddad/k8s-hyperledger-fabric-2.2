package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (rc *ResourceTypesContract) SetPrivateData(
	ctx contractapi.TransactionContextInterface,
	id string,
	privateName string,
) error {
	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return err
	}

	if err = ctx.GetStub().PutPrivateData(
		mspID+"ResourceTypesPrivateData",
		id,
		[]byte(privateName),
	); err != nil {
		return err
	}

	return nil
}
