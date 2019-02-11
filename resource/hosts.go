package resource

import (
	"encoding/json"
	"errors"
	"fmt"
	"intel/isecl/threat-detection-service/repository"
	"intel/isecl/threat-detection-service/types"
	"net/http"

	"github.com/gorilla/mux"
)

func SetHosts(r *mux.Router, db repository.TDSDatabase) {
	r.HandleFunc("/hosts", coerceJSON(errorHandler(createHost(db)))).Methods("POST")
	r.HandleFunc("/hosts", errorHandler(nil)).Methods("GET")
	r.HandleFunc("/hosts/{id}", errorHandler(nil)).Methods("DELETE")
	r.HandleFunc("/hosts/{id}", errorHandler(nil)).Methods("GET")
	r.HandleFunc("/hosts/{id}", coerceJSON(errorHandler(nil))).Methods("PATCH")
}

func coerceJSON(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte(`Content-Type must be "application/json"`))
		} else {
			h(w, r)
		}
	}
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
		//return db.HostRepository().Create(host)
		return nil
	}
}

func getHost(db repository.TDSDatabase) errorHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		id := mux.Vars(r)["id"]
		fmt.Println(id)
		// h := db.HostRepository().Retrieve(id)
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
