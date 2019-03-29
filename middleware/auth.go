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
			// the username from r.BasicAuth() is either username or host ID on TDAgent
			username, password, ok := r.BasicAuth()
			log.Debug("Attempting to authenticate user: ", username)
			if ok {
				// fetch by user
				user, err := getUserByNameOrId(username, u)
				if err != nil {
					log.WithError(err).Error("BasicAuth failure: could not retrieve user")
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				if err := user.CheckPassword([]byte(password)); err != nil {
					log.WithError(err).Error("BasicAuth failure: password mismatch, user", username)
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				roles, err_role := getRolesByNameOrID(username, u)
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

func getUserByNameOrId(userNameOrId string, u repository.UserRepository) (*types.User, error) {
	user, err := u.Retrieve(types.User{Name: userNameOrId})
	if err != nil {
		user, err = u.Retrieve(types.User{ID: userNameOrId})
		if err != nil {
			return nil, err
		}
	}
	return user, nil
}

func getRolesByNameOrID(userNameOrId string, u repository.UserRepository) (userRoles []types.Role, err error) {
	roles, err_role := u.GetRoles(types.User{Name: userNameOrId})
	if err_role != nil || len(roles) == 0 {
		roles, err_role = u.GetRoles(types.User{ID: userNameOrId})
		if err != nil {
			return nil, err_role
		}
	}
	return roles, nil
}
