package models

import (
	"admin-api/internal/k8client"
	"admin-api/types"

	"xorm.io/xorm"
)

type CertificateAuthorities interface {
	Create(*types.OrgCertificateAuthority) error
}

func NewCertificateAuthorities(db *xorm.Engine, k8c k8client.Client) CertificateAuthorities {
	return &certificateAuthorities{db, k8c}
}

type certificateAuthorities struct {
	db  *xorm.Engine
	k8c k8client.Client
}

func (c *certificateAuthorities) Create(ca *types.OrgCertificateAuthority) (err error) {
	if err = types.Validate(ca); err != nil {
		return err
	}

	return nil
}
