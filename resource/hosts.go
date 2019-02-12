package resource

import (
	"encoding/json"
	"errors"
	"fmt"
	"intel/isecl/threat-detection-service/repository"
	"intel/isecl/threat-detection-service/types"
	"net/http"

	"github.com/gorilla/handlers"

	"github.com/gorilla/mux"
)

func SetHosts(r *mux.Router, db repository.TDSDatabase) {
	r.Handle("/hosts", handlers.ContentTypeHandler(createHost(db), "application/json"))
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
		err = db.HostRepository().Create(h)
		if err != nil {
			return err
		}
		w.WriteHeader(http.StatusCreated) // HTTP 201
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(&h)
		if err != nil {
			return err
		}
		return nil
	}
}

func getHost(db repository.TDSDatabase) errorHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		id := mux.Vars(r)["id"]
		fmt.Println(id)
		h, err := db.HostRepository().Retrieve(types.Host{ID: id})
		if err != nil {
			return err
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(&h)
		if err != nil {
			return err
		}
		return nil
	}
}

func deleteHost(db repository.TDSDatabase) errorHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		return nil
	}
}

func queryHosts(db repository.TDSDatabase) errorHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		return nil
	}
}

func updateHost(db repository.TDSDatabase) errorHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		return nil
	}
}
