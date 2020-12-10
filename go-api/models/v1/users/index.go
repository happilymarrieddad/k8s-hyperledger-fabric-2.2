package users

import "k8s-hyperledger-fabric-2.2/go-api/models"

func Index() (users *models.Users, err error) {
	users = &mockUsers

	return
}
