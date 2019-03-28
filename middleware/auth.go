package middleware

import (
	"intel/isecl/tdservice/repository"
	"intel/isecl/tdservice/types"
	"net/http"

	"github.com/gorilla/mux"
	 _"github.com/gorilla/context"
	log "github.com/sirupsen/logrus"
)

func NewBasicAuth(u repository.UserRepository) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, password, ok := r.BasicAuth()
			log.Debug("Attempting to authenticate user: ", username)
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
					log.WithError(err).Error("BasicAuth failure: password mismatch, user", username)
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				next.ServeHTTP(w, r)
			} else {
				log.Info("No Basic Auth provided")
				w.WriteHeader(http.StatusUnauthorized)
			}
		})
	}
}
