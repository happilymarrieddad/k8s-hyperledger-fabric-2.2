package resourcetypes

import (
	"encoding/json"
	"net/http"

	"k8s-hyperledger-fabric-2.2/go-api/hyperledger"
	ResourceTypesModel "k8s-hyperledger-fabric-2.2/go-api/models/v1/resourcetypes"
)

func Index(clients *hyperledger.Clients) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		resourcetypes, err := ResourceTypesModel.Index(clients)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		packet, err := json.Marshal(resourcetypes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Write(packet)
	}
}
