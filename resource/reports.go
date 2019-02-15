package resource

import (
	"encoding/json"
	"intel/isecl/threat-detection-service/repository"
	"intel/isecl/threat-detection-service/types"
	"net/http"

	"github.com/gorilla/handlers"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func SetReports(r *mux.Router, db repository.TDSDatabase) {
	r.Handle("/reports", handlers.ContentTypeHandler(createReport(db), "application/json")).Methods("POST")
	r.Handle("/reports", queryReport(db)).Methods("GET")
	r.Handle("/reports/{id}", getReport(db)).Methods("GET")
}

func createReport(db repository.TDSDatabase) errorHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var report types.Report
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&report)
		if err != nil {
			log.WithError(err).Error("failed to decode input body as types.Report")
			return err
		}
		db.ReportRepository().Create(&report)

		w.WriteHeader(http.StatusCreated) // HTTP 201
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		err = encoder.Encode(report)
		if err != nil {
			log.WithError(err).Error("failed to encode response body as types.Report")
			return err
		}
		return nil
	}
}

func queryReport(db repository.TDSDatabase) errorHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		return nil
	}
}

func getReport(db repository.TDSDatabase) errorHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		return nil
	}
}
