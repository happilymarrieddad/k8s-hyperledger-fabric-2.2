package resources

import (
	"encoding/json"
	"net/http"

	"k8s-hyperledger-fabric-2.2/go-api/models"

	"k8s-hyperledger-fabric-2.2/go-api/hyperledger"
	ResourcesModel "k8s-hyperledger-fabric-2.2/go-api/models/v1/resources"
)

func Store(clients *hyperledger.Clients) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var resource models.Resource
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		err := decoder.Decode(&resource)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		newResource, err := ResourcesModel.Store(
			clients,
			resource.Name,
			resource.ResourceTypeID,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		packet, err := json.Marshal(newResource)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(packet)
	}
}
