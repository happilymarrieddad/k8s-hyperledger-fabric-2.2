package types

type OrgCertificateAuthority struct {
	ID      int64  `xorm:"'id' pk autoincr"`
	OrgID   int64  `validate:"required" xorm:"org_id"`
	Name    string `validate:"required" xorm:"name"`
	DNSName string `validate:"required" xorm:"dns_name"`
	IsRoot  bool   `xorm:"is_root"`
}

func (o *OrgCertificateAuthority) TableName() string {
	return `org_certificate_authorities`
}
