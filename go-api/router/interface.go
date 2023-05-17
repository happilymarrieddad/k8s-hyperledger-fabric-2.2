package router

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"k8s-hyperledger-fabric-2.2/go-api/hyperledger"
	V1Router "k8s-hyperledger-fabric-2.2/go-api/routes/v1"
)

const (
	staticDir = "/static/"
)

type Service interface {
	GetRawRouter() *mux.Router
}

func GetRouter() Service {
	r := Router{
		RawRouter: mux.NewRouter().StrictSlash(true),
	}

	configPath := os.Getenv("HYPERLEDGER_CONFIG_PATH")
	MSPID := os.Getenv("HYPERLEDGER_MSP_ID")
	if len(MSPID) == 0 {
		MSPID = "ibm"
	}

	if len(configPath) == 0 {
		panic("ENV var 'HYPERLEDGER_CONFIG_PATH' is not set. unable to connect to network")
	}

	clients := hyperledger.NewClientMap(
		"test-network",
		configPath,
	)

	_, err := clients.AddClient(
		"Admin",
		MSPID,
		"mainchannel",
	)
	if err != nil {
		panic(err)
	}

	if err = clients.SetupChannelListener(MSPID, "mainchannel", "resources"); err != nil {
		panic(err)
	}

	if err = clients.SetupChannelListener(MSPID, "mainchannel", "resource_types"); err != nil {
		panic(err)
	}

	r.RawRouter.
		PathPrefix(staticDir).
		Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir("."+staticDir))))

	for _, route := range GetRoutes() {
		r.RawRouter.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	for name, pack := range V1Router.GetRoutes(clients) {
		r.AttachSubRouterWithMiddleware(name, pack.Routes, pack.Middleware)
	}

	return r
}
