package users

import (
	"k8s-hyperledger-fabric-2.2/go-api/models"
)

var mockUsers models.Users

func init() {
	usr, _ := models.NewUser("Nick", "Kotenberg", "nick@mail.com", "1234")

	mockUsers = models.Users{
		*usr,
	}
}
