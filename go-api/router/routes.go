package router

import (
	"k8s-hyperledger-fabric-2.2/go-api/models"
	HomeHandler "k8s-hyperledger-fabric-2.2/go-api/routes/home"
	StatusHandler "k8s-hyperledger-fabric-2.2/go-api/routes/status"
)

func GetRoutes() models.Routes {
	return models.Routes{
		models.Route{Name: "Home", Method: "GET", Pattern: "/", HandlerFunc: HomeHandler.Index},
		models.Route{Name: "Status", Method: "GET", Pattern: "/status", HandlerFunc: StatusHandler.Index},
	}
}
