package resource

import (
	"fmt"
	"intel/isecl/tdservice/repository"
	"intel/isecl/tdservice/version"
	"net/http"

	"github.com/gorilla/mux"
)

func SetVersion(r *mux.Router, db repository.TDSDatabase) {
	r.Handle("/version", getVersion()).Methods("GET")
}

func getVersion() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		verStr := fmt.Sprintf("%s-%s", version.Version, version.GitHash)
		w.Write([]byte(verStr))
	})
}
