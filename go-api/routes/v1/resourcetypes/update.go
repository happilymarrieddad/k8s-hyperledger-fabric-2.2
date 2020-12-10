package resourcetypes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"k8s-hyperledger-fabric-2.2/go-api/models"

	ResourceTypesModel "k8s-hyperledger-fabric-2.2/go-api/models/v1/resourcetypes"
)

func Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var opts ResourceTypesModel.UpdateOpts
		var resourcetype models.ResourceType
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		if err := decoder.Decode(&resourcetype); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if r.Method == "PUT" {
			opts.Replace = true
		}

		updatedResourceType, err := ResourceTypesModel.Update(id, &resourcetype, &opts)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		packet, err := json.Marshal(updatedResourceType)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(packet)
	}
}
