package repository

import "github.com/1sh-repalto/url-monitoring-api/internal/model"

type UserRepository interface {
	GetUserByEmail (email string) (*model.User, error)
	CreateUser (user *model.User) error
}