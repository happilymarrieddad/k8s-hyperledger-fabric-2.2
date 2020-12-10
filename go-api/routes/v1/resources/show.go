package resources

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	ResourcesModel "k8s-hyperledger-fabric-2.2/go-api/models/v1/resources"
	"k8s-hyperledger-fabric-2.2/go-api/hyperledger"
)

func Show(clients *hyperledger.Clients) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		id := vars["id"]

		rawResource, err := ResourcesModel.Show(clients, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		packet, err := json.Marshal(rawResource)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(packet)
	}
}
