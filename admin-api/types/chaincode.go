package types

type chaincode string

const (
	// TODO: add real chaincodes at some point
	OrdersChaincode chaincode = "orders"
)

type OrgPeerChaincode struct {
	ID                    int64     `xorm:"'id' pk autoincr"`
	OrgID                 int64     `xorm:"org_id"`
	PeerID                int64     `xorm:"peer_id"`
	Name                  chaincode `xorm:"name"`
	CurrentSequenceNumber int64     `xorm:"current_sequence_number"`
	CurrentVersionNumber  int64     `xorm:"current_version_number"`
}

func (o *OrgPeerChaincode) TableName() string {
	return `org_peer_chaincodes`
}
