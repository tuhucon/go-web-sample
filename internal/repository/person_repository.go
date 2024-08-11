package repository

import "fileserver/internal/model"

type PersonRepository interface {
	FindAll() ([]*model.Person, error)
}
