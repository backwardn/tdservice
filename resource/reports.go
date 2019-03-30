/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package resource

import (
	"encoding/json"
	"errors"
	"intel/isecl/tdservice/constants"
	"intel/isecl/tdservice/context"
	"intel/isecl/tdservice/repository"
	"intel/isecl/tdservice/types"
	"net/http"
	"time"
	"bytes"
	"strings"
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
		reportsBuffer := new(bytes.Buffer)
		reportsBuffer.ReadFrom(r.Body)
		dec := json.NewDecoder(strings.NewReader(reportsBuffer.String()))
		dec.DisallowUnknownFields()
		err := dec.Decode(&report)
		if err != nil {
			log.WithError(err).Error("failed to decode input body as types.Report")
			return err
		}
		if report.HostID == "" {
			log.Error("report is not associated with a HostID")
			return errors.New("report is not associated with a HostID")
		}

		// Check query authority
		roles := context.GetUserRoles(r)
		actionAllowed := false
		for _, role := range roles {
			if role.Name == constants.HostSelfUpdateGroupName && role.Domain == report.HostID {
				actionAllowed = true
				break
			}
		}
		if !actionAllowed {
			return &privilegeError{Message: "privilege error: create report",
				StatusCode: http.StatusForbidden}
		}

		log.WithField("report", report).Info("creating report")
		created, err := db.ReportRepository().Create(report)

		w.WriteHeader(http.StatusCreated) // HTTP 201
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(&created)
		if err != nil {
			log.WithError(err).Error("failed to encode response body as types.Report")
			return err
		}
		return nil
	}
}

func queryReport(db repository.TDSDatabase) errorHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {

		// Check query authority
		roles := context.GetUserRoles(r)
		actionAllowed := false
		for _, role := range roles {
			if role.Name == constants.AdminGroupName {
				actionAllowed = true
				break
			}
		}
		if !actionAllowed {
			return &privilegeError{Message: "privilege error: query report",
				StatusCode: http.StatusForbidden}
		}

		q := r.URL.Query()
		hostname := q.Get("hostname")
		hostID := q.Get("hostid")
		fromStr := q.Get("from")
		toStr := q.Get("to")
		var from time.Time
		var to time.Time
		var err error
		if fromStr != "" {
			from, err = time.Parse(time.RFC3339, fromStr)
			if err != nil {
				log.WithError(err).WithField("from", fromStr).Error("failed to parse RFC3339 date")
				// explicitly return bad request
				return &resourceError{Message: err.Error(), StatusCode: http.StatusBadRequest}
			}
		}
		if toStr != "" {
			to, err = time.Parse(time.RFC3339, toStr)
			if err != nil {
				log.WithError(err).WithField("to", toStr).Error("failed to parse RFC3339 ")
				// explicitly return bad request
				return &resourceError{Message: err.Error(), StatusCode: http.StatusBadRequest}
			}
		}
		filter := repository.ReportFilter{}
		filter.Hostname = hostname
		filter.HostID = hostID
		filter.From = from
		filter.To = to

		reports, err := db.ReportRepository().RetrieveByFilterCriteria(filter)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&reports)
		return nil
	}
}

func getReport(db repository.TDSDatabase) errorHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		id := mux.Vars(r)["id"]

		// Check query authority
		roles := context.GetUserRoles(r)
		actionAllowed := false
		for _, role := range roles {
			if role.Name == constants.AdminGroupName {
				actionAllowed = true
				break
			}
		}
		if !actionAllowed {
			return &privilegeError{Message: "privilege error: get report",
				StatusCode: http.StatusForbidden}
		}

		h, err := db.ReportRepository().Retrieve(types.Report{ID: id})
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
