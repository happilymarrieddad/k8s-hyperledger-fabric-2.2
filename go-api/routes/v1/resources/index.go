package resources

import (
	"encoding/json"
	"net/http"

	ResourcesModel "k8s-hyperledger-fabric-2.2/go-api/models/v1/resources"
	"k8s-hyperledger-fabric-2.2/go-api/hyperledger"
)

func Index(clients *hyperledger.Clients) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		resources, err := ResourcesModel.Index(clients)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		packet, err := json.Marshal(resources)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Write(packet)
	}
}
