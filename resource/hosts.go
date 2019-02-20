package resource

import (
	"encoding/json"
	"errors"
	"intel/isecl/threat-detection-service/repository"
	"intel/isecl/threat-detection-service/types"
	"net/http"

	"github.com/gorilla/handlers"
	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

func SetHosts(r *mux.Router, db repository.TDSDatabase) {
	r.Handle("/hosts", handlers.ContentTypeHandler(createHost(db), "application/json")).Methods("POST")
	r.Handle("/hosts", queryHosts(db)).Methods("GET")
	r.Handle("/hosts/{id}", deleteHost(db)).Methods("DELETE")
	r.Handle("/hosts/{id}", getHost(db)).Methods("GET")
	r.Handle("/hosts/{id}", handlers.ContentTypeHandler(updateHost(db), "application/json")).Methods("PATCH")
}

func createHost(db repository.TDSDatabase) errorHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var h types.Host
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		err := dec.Decode(&h.HostInfo)
		if err != nil {
			return err
		}
		// validate host
		if h.Hostname == "" {
			return errors.New("hostname is empty")
		}
		if h.OS == "" {
			return errors.New("os is empty")
		}
		if h.Version == "" {
			return errors.New("version is empty")
		}
		if err != nil {
			return err
		}
		// add the host
		// the server, on a separate go routine, will periodically ping all registered hosts to update their status
		created, err := db.HostRepository().Create(h)
		if err != nil {
			return err
		}
		w.WriteHeader(http.StatusCreated) // HTTP 201
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(&created)
		if err != nil {
			return err
		}
		return nil
	}
}

func getHost(db repository.TDSDatabase) errorHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		id := mux.Vars(r)["id"]
		h, err := db.HostRepository().Retrieve(types.Host{ID: id})
		if err != nil {
			log.WithError(err).WithField("id", id).Info("failed to retrieve host")
			return err
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(&h)
		if err != nil {
			log.WithError(err).Error("failed to encode json response")
			return err
		}
		return nil
	}
}

func deleteHost(db repository.TDSDatabase) errorHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		id := mux.Vars(r)["id"]
		if err := db.HostRepository().Delete(types.Host{ID: id}); err != nil {
			log.WithError(err).WithField("id", id).Info("failed to delete host")
			return err
		}
		w.WriteHeader(http.StatusNoContent)
		return nil
	}
}

func queryHosts(db repository.TDSDatabase) errorHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		// check for query parameters
		log.WithField("query", r.URL.Query()).Trace("query hosts")
		hostname := r.URL.Query().Get("hostname")
		version := r.URL.Query().Get("version")
		build := r.URL.Query().Get("build")
		os := r.URL.Query().Get("os")
		status := r.URL.Query().Get("status")

		filter := types.Host{
			HostInfo: types.HostInfo{
				Hostname: hostname,
				Version:  version,
				Build:    build,
				OS:       os,
			},
			Status: status,
		}

		hosts, err := db.HostRepository().RetrieveAll(filter)
		if err != nil {
			log.WithError(err).WithField("filter", filter).Info("failed to retrieve hosts")
			return err
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(hosts)
		return nil
	}
}

func updateHost(db repository.TDSDatabase) errorHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		return nil
	}
}
