package middleware

import (
	"intel/isecl/threat-detection-service/repository"
	"intel/isecl/threat-detection-service/types"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func NewBasicAuth(u repository.UserRepository) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, password, ok := r.BasicAuth()
			if ok {
				// fetch by user
				user, err := u.Retrieve(types.User{Name: username})
				if err != nil {
					// log this
					log.WithError(err).Error("BasicAuth failure: could not retrieve user")
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				if err := user.CheckPassword([]byte(password)); err != nil {
					log.WithError(err).Error("BasicAuth failure: password mismatch")
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				next.ServeHTTP(w, r)
			}
		})
	}
}
