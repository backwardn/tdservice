package repository

import (
	"github.com/google/uuid"
)

// UUID returns a random uuid string and error if there is any
func UUID() (string, error) {

	uuid, err := uuid.NewRandom()
	if err == nil {
		return uuid.String(), nil
	} else {
		return "", err
	}
}
