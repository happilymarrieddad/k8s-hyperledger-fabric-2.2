package types

type OrgPeer struct {
	ID    int64  `xorm:"'id' pk autoincr"`
	OrgID int64  `validate:"required" xorm:"org_id"`
	Name  string `validate:"required" xorm:"name"`
	// Need to autoincrement this at some point based on the current peers
	Number  int64  `validate:"required" xorm:"Number"`
	DNSName string `validate:"required,alpha,min=4,max=10" xorm:"dns_name"`
	Created bool   `xorm:"created"`
}

func (o *OrgPeer) TableName() string {
	return `org_peers`
}
