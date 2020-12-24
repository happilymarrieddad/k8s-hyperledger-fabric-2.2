package main

// https://hyperledger-fabric.readthedocs.io/en/latest/chaincode4ade.html
import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	cc, err := contractapi.NewChaincode(&ResourceTypesContract{})

	if err != nil {
		panic(err.Error())
	}

	if err := cc.Start(); err != nil {
		panic(err.Error())
	}
}

// ResourceTypesContract contract for handling writing and reading from the world state
type ResourceTypesContract struct {
	contractapi.Contract
}

// InitLedger adds a base set of resource types to the ledger
func (rc *ResourceTypesContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	return nil
}
