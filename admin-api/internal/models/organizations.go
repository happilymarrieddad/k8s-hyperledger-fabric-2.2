package models

import (
	"admin-api/internal/k8client"
	"admin-api/types"
	"fmt"

	"github.com/pkg/errors"
	"xorm.io/xorm"
)

type Organizations interface {
	Create(*types.Organization) (storageDeploymentPodID string, err error)
	FindByID(id int64) (*types.Organization, error)
	Delete(id int64) error
}

func NewOrganizations(db *xorm.Engine, k8c k8client.Client) Organizations {
	return &organizations{
		db, k8c,
		// Models
	}
}

type organizations struct {
	db  *xorm.Engine
	k8c k8client.Client
	// models

}

func (m *organizations) Create(org *types.Organization) (storageDeploymentPodID string, err error) {
	if err = types.Validate(org); err != nil {
		return
	}

	sesh := m.db.NewSession()
	defer sesh.Close()

	k8sNetworkCleanup := func() {
		if errc := m.k8c.DeleteNamespace(org.NetworkName); errc != nil {
			fmt.Println("unable to cleanup namespace for org ", org.Name)
		}
	}

	if err = sesh.Begin(); err != nil {
		return
	}

	if _, err = sesh.Insert(org); err != nil {
		sesh.Rollback()
		return "", errors.WithMessage(err, "unable to insert org")
	}

	nmsp, err := m.k8c.CreateNamespace(org.NetworkName)
	if err != nil {
		sesh.Rollback()
		return "", errors.WithMessage(err, "unable to create namespace")
	}

	if _, err := sesh.ID(org.ID).Cols("namespace_created").Update(&types.Organization{NamespaceCreated: true}); err != nil {
		k8sNetworkCleanup()
		sesh.Rollback()
		return "", errors.WithMessage(err, "unable to update 'namespace_created'")
	}

	// Set it on the org so that it will be available on the return
	org.NamespaceCreated = true

	// Create a storage
	stgName := fmt.Sprintf("%s-storage", org.NetworkName)
	if _, err := m.k8c.CreateNamespaceStorage(stgName, nmsp.Name); err != nil {
		k8sNetworkCleanup()
		sesh.Rollback()
		return "", errors.WithMessage(err, "unable to create storage")
	}

	// Create an Org deployment
	storageDeploymentPodID, err = m.k8c.CreateNamespaceStorageDeployment(stgName, nmsp.Name)
	if err != nil {
		k8sNetworkCleanup()
		sesh.Rollback()
		return "", errors.WithMessage(err, "unable to create storage deployment")
	}

	// Copy required files to the pod in the deployment

	// Create CA(s)

	// Create certs using ca(s)

	if err = sesh.Commit(); err != nil {
		k8sNetworkCleanup()
		err = errors.WithMessage(err, "unable to commit organization")
	}
	return
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
