package resource

import (
	"encoding/json"
	"errors"
	"fmt"
	"intel/isecl/lib/common/validation"
	"intel/isecl/lib/common/crypt"
	"intel/isecl/tdservice/repository"
	"intel/isecl/tdservice/types"
	consts "intel/isecl/tdservice/constants"

	"net/http"

	"github.com/gorilla/handlers"
	log "github.com/sirupsen/logrus"
	
	"golang.org/x/crypto/bcrypt"
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
		var valid_err error
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		err := dec.Decode(&h.HostInfo)
		if err != nil {
			return err
		}
		// validate host
		if h.Hostname == "" {
			return errors.New("hostname is invalid")
		}
		valid_err = validation.ValidateHostname(h.Hostname)
		if valid_err != nil {
			return fmt.Errorf("hostname validation fail: %s", valid_err.Error())
		}
		// validate os
		if h.OS != "linux" && h.OS != "windows" {
			return errors.New("os is invalid")
		}
		// validate version
		if h.Version == "" || len(h.Version) > 32 {
			return errors.New("version is invalid")
		}

		// valid_err = validation.ValidateRestrictedString(h.Version, "a-z\\-.0-9")
		// if valid_err != nil {
		// 	return fmt.Errorf("validation fail: %s", valid_err.Error())
		// }
		// validate build
		if h.Build == "" || len(h.Build) > 32 {
			return errors.New("build is invalid")
		}
		// valid_err = validation.ValidateRestrictedString(h.Build, "a-z0-9")
		// if valid_err != nil {
		// 	return fmt.Errorf("validation fail: %s", valid_err.Error())
		// }

		if err != nil {
			return err
		}
		// add the host
		// the server, on a separate go routine, will periodically ping all registered hosts to update their status, for now, assume online
		h.Status = "Reserve for future implementation"
		created, err := db.HostRepository().Create(h)
		if err != nil {
			return err
		}
		// create the user and roles that represents the new domain. API endpoints that are restricted to updates only from the newly created 
		// hosts shall be protected with the role. 
		rand, err := crypt.GetRandomBytes(consts.PasswordRandomLength)
		if err != nil {
			return err
		}
		hash, err := bcrypt.GenerateFromPassword(rand, bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		host_user_role := types.Role{Name: consts.HostSelfUpdateGroupName,
									 Domain: created.ID}
		host_user := types.User{PasswordHash: hash,
								Roles : []types.Role{host_user_role}}

		user, err := db.UserRepository().Create(host_user)
		if err != nil {
			return err
		}
		resp  := types.HostCreateResponse{}
		resp.HostInfo = created.HostInfo
		resp.User = user.ID
		resp.Token = rand

		w.WriteHeader(http.StatusCreated) // HTTP 201
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(&resp)
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
		// status := r.URL.Query().Get("status")

		filter := types.Host{
			HostInfo: types.HostInfo{
				Hostname: hostname,
				Version:  version,
				Build:    build,
				OS:       os,
			},
			Status: "Reserve for future implementation",
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
