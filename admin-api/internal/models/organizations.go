package models

import (
	"admin-api/internal/k8client"
	"admin-api/types"
	"fmt"

	"xorm.io/xorm"
)

type Organizations interface {
	Create(*types.Organization) error
	FindByID(id int64) (*types.Organization, error)
	Delete(id int64) error
}

func NewOrganizations(db *xorm.Engine, k8c k8client.Client) Organizations {
	return &organizations{db, k8c}
}

type organizations struct {
	db  *xorm.Engine
	k8c k8client.Client
}

func (m *organizations) Create(org *types.Organization) error {
	if err := types.Validate(org); err != nil {
		return err
	}

	sesh := m.db.NewSession()
	defer sesh.Close()

	if err := sesh.Begin(); err != nil {
		return err
	}

	if _, err := sesh.Insert(org); err != nil {
		sesh.Rollback()
		return err
	}

	if _, err := m.k8c.CreateNamespace(org.NetworkName); err != nil {
		sesh.Rollback()
		return err
	}

	if _, err := sesh.ID(org.ID).Cols("namespace_created").Update(&types.Organization{NamespaceCreated: true}); err != nil {
		_ = m.k8c.DeleteNamespace(org.NetworkName)
		sesh.Rollback()
		return err
	}

	// Set it on the org so that it will be available on the return
	org.NamespaceCreated = true

	// Create org resources
	// CA(s)
	// Certs
	// TODO: decide if done here or in the UI
	// Peer(s)
	// Chaincode(s)

	return sesh.Commit()
}

func (m *organizations) FindByID(id int64) (*types.Organization, error) {
	org := new(types.Organization)

	if has, err := m.db.ID(id).Where("active = ?", true).Get(org); err != nil {
		return nil, err
	} else if !has {
		return nil, fmt.Errorf("unable to find organization: %d", id)
	}

	return org, nil
}

func (m *organizations) Delete(id int64) error {
	sesh := m.db.NewSession()
	defer sesh.Close()

	if err := sesh.Begin(); err != nil {
		return err
	}

	org, err := m.FindByID(id)
	if err != nil {
		sesh.Rollback()
		return err
	}

	if _, err := sesh.ID(id).Cols("active").Update(&types.Organization{Active: false}); err != nil {
		sesh.Rollback()
		return err
	}

	if _, err := sesh.ID(org.ID).Cols("namespace_created").Update(&types.Organization{NamespaceCreated: false}); err != nil {
		sesh.Rollback()
		return err
	}

	// This probably should be before the update to the namespace created but if this happens it's a major issue...
	// ALL reasources will be removed including all data for this org... very dangerous
	// TODO: decide if this is even something we should do
	if err := m.k8c.DeleteNamespace(org.Name); err != nil {
		sesh.Rollback()
		return err
	}

	return sesh.Commit()
}
