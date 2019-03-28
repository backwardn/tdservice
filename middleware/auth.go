package middleware

import (
	"intel/isecl/tdservice/context"
	"intel/isecl/tdservice/repository"
	"intel/isecl/tdservice/types"
	"net/http"

	_ "github.com/gorilla/context"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func NewBasicAuth(u repository.UserRepository) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userNameorId, password, ok := r.BasicAuth()
			log.Debug("Attempting to authenticate user: ", userNameorId)
			if ok {
				// fetch by user
				user, err := u.Retrieve(types.User{Name: userNameorId})
				if err != nil {
					user, err = u.Retrieve(types.User{ID: userNameorId})
					if err != nil {
						log.WithError(err).Error("BasicAuth failure: could not retrieve user")
						w.WriteHeader(http.StatusUnauthorized)
						return
					}
				}
				if err := user.CheckPassword([]byte(password)); err != nil {
					log.WithError(err).Error("BasicAuth failure: password mismatch, user", userNameorId)
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				roles, err_role := u.GetRoles(types.User{Name: username})
				if err_role != nil {
					log.WithError(err).Error("Database error: unable to retrive roles")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				context.SetUserRoles(r, roles)
				next.ServeHTTP(w, r)
			} else {
				log.Info("No Basic Auth provided")
				w.WriteHeader(http.StatusUnauthorized)
			}
		})
	}
}
