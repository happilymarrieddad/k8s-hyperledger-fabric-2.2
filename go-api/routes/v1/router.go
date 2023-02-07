package v1

import (
	"net/http"

	"k8s-hyperledger-fabric-2.2/go-api/hyperledger"
	"k8s-hyperledger-fabric-2.2/go-api/models"
	ResourcesHandler "k8s-hyperledger-fabric-2.2/go-api/routes/v1/resources"
	ResourceTypesHandler "k8s-hyperledger-fabric-2.2/go-api/routes/v1/resourcetypes"
	UsersHandler "k8s-hyperledger-fabric-2.2/go-api/routes/v1/users"
)

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}

func GetRoutes(clients *hyperledger.Clients) map[string]models.SubRoutePackage {
	return map[string]models.SubRoutePackage{
		"/v1": {
			Middleware: Middleware(),
			Routes: models.Routes{
				// Users
				models.Route{Name: "UsersIndex", Method: "GET", Pattern: "/users", HandlerFunc: UsersHandler.Index()},
				models.Route{Name: "UsersStore", Method: "POST", Pattern: "/users", HandlerFunc: UsersHandler.Store()},
				models.Route{Name: "UsersReplace", Method: "PUT", Pattern: "/users/{id}", HandlerFunc: UsersHandler.Update()},
				models.Route{Name: "UsersUpdate", Method: "PATCH", Pattern: "/users/{id}", HandlerFunc: UsersHandler.Update()},
				models.Route{Name: "UsersDestroy", Method: "DELETE", Pattern: "/users/{id}", HandlerFunc: UsersHandler.Destroy()},
				// ResourceTypes
				models.Route{Name: "ResourceTypesIndex", Method: "GET", Pattern: "/resource_types", HandlerFunc: ResourceTypesHandler.Index(clients)},
				models.Route{Name: "ResourceTypesStore", Method: "POST", Pattern: "/resource_types", HandlerFunc: ResourceTypesHandler.Store()},
				models.Route{Name: "ResourceTypesReplace", Method: "PUT", Pattern: "/resource_types/{id}", HandlerFunc: ResourceTypesHandler.Update()},
				models.Route{Name: "ResourceTypesUpdate", Method: "PATCH", Pattern: "/resource_types/{id}", HandlerFunc: ResourceTypesHandler.Update()},
				models.Route{Name: "ResourceTypesDestroy", Method: "DELETE", Pattern: "/resource_types/{id}", HandlerFunc: ResourceTypesHandler.Destroy()},
				// Resources
				models.Route{Name: "ResourcesIndex", Method: "GET", Pattern: "/resources", HandlerFunc: ResourcesHandler.Index(clients)},
				models.Route{Name: "ResourcesStore", Method: "POST", Pattern: "/resources", HandlerFunc: ResourcesHandler.Store(clients)},
				models.Route{Name: "ResourcesReplace", Method: "PUT", Pattern: "/resources/{id}", HandlerFunc: ResourcesHandler.Update(clients)},
				models.Route{Name: "ResourcesUpdate", Method: "PATCH", Pattern: "/resources/{id}", HandlerFunc: ResourcesHandler.Update(clients)},
				models.Route{Name: "ResourcesDestroy", Method: "DELETE", Pattern: "/resources/{id}", HandlerFunc: ResourcesHandler.Destroy(clients)},
				models.Route{Name: "ResourcesShow", Method: "GET", Pattern: "/resources/{id}", HandlerFunc: ResourcesHandler.Show(clients)},
			},
		},
	}
}
