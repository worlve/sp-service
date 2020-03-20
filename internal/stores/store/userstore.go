package store

import "github.com/worlve/sp-service/internal/models/appuser"

// UserStore defines the required functionality for any associated store.
type UserStore interface {
	GetUser(userGUID string) (appuser.User, error)
}
