package resourcetypes

import (
	"encoding/json"
	"net/http"

	ResourceTypesModel "k8s-hyperledger-fabric-2.2/go-api/models/v1/resourcetypes"
)

func Index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		resourcetypes, err := ResourceTypesModel.Index()
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
