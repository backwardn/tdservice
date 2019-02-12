package resource

import (
	"intel/isecl/threat-detection-service/repository"
	"net/http"

	"github.com/gorilla/handlers"

	"github.com/gorilla/mux"
)

func SetReports(r *mux.Router, db repository.TDSDatabase) {
	r.Handle("/reports", handlers.ContentTypeHandler(createReport(db), "application/json")).Methods("POST")
	r.Handle("/reports", queryReport(db)).Methods("GET")
	r.Handle("/reports/{id}", getReport(db)).Methods("GET")
}

func createReport(db repository.TDSDatabase) errorHandlerFunc {
	return errorHandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		return nil
	})
}

func queryReport(db repository.TDSDatabase) errorHandlerFunc {
	return errorHandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		return nil
	})
}

func getReport(db repository.TDSDatabase) errorHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		return nil
	}
}
