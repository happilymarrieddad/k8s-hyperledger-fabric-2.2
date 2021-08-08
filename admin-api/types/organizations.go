package types

type Organization struct {
	ID               int64  `xorm:"'id' pk autoincr"`
	Name             string `validate:"required" xorm:"name"`
	NetworkName      string `validate:"required,alpha,min=4,max=10" xorm:"network_name"`
	NamespaceCreated bool   `xorm:"namespace_created"`
	Active           bool   `xorm:"active"`
	// Network Resource Information
	NumberOfCAs   int64 `validate:"required,min=1,max=3" xorm:"number_of_cas"`
	NumberOfPeers int64 `validate:"required,min=2,max=10" xorm:"number_of_peers"`
	// Network Flags
	NetworkActive bool `xorm:"network_active"`
}

func (o *Organization) TableName() string {
	return `organizations`
}
