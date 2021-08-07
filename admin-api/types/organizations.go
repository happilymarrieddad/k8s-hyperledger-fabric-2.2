package types

type Organization struct {
	ID               int64  `xorm:"'id' pk autoincr"`
	Name             string `validate:"required" xorm:"name"`
	NetworkName      string `validate:"required,alpha,min=4,max=10" xorm:"network_name"`
	NamespaceCreated bool   `xorm:"namespace_created"`
	Active           bool   `xorm:"active"`
}

func (n *Organization) TableName() string {
	return `organizations`
}
