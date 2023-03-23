package mgin

import "github.com/pkg/errors"

type User interface {
	GetID() uint
	GetUsername() string
	GetPassword() string
}

var ErrUserIDNotFound = errors.New("user not found")
var ErrUsernameNotFound = errors.New("username not found")

type UserStore interface {
	Find(userID uint) (User, error)
	FindByUsername(username string) (User, error)
	Save(user User) error
}

type SettingStore interface {
	Get(userID uint, key string) (string, error)
	Set(userID uint, key string, value string) error
}
