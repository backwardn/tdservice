package resource

import (
	"intel/isecl/threat-detection-service/repository"
	"net/http"

	"github.com/gorilla/mux"
)

func SetReports(r *mux.Router, db repository.TDSDatabase) {
	r.HandleFunc("/reports", coerceJSON(errorHandler(createHost(db)))).Methods("POST")
	r.HandleFunc("/reports", errorHandler(nil)).Methods("GET")
	r.HandleFunc("/reports/{id}", errorHandler(nil)).Methods("GET")
}

func createReport(db repository.TDSDatabase) errorHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		return nil
	}
}
