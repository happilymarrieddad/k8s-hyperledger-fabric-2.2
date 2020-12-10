package resourcetypes

import (
	"net/http"

	"github.com/gorilla/mux"

	ResourceTypesModel "k8s-hyperledger-fabric-2.2/go-api/models/v1/resourcetypes"
)

func Destroy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		id := vars["id"]

		if err := ResourceTypesModel.Destroy(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte("Success"))
	}
}
