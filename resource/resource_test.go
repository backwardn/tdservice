package resource

import (
	"intel/isecl/tdservice/repository"

	"github.com/gorilla/mux"
)

func setupRouter(db repository.TDSDatabase) *mux.Router {
	r := mux.NewRouter().PathPrefix("/tds").Subrouter()
	SetHosts(r, db)
	SetReports(r, db)
	return r
}
