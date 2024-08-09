package controller

import (
	"fileserver/internal/model"
	"fileserver/internal/repository"
	"net/http"
)

type PersonController struct {
	BaseController
	PersonRepository repository.PersonRepository
}

func (c *PersonController) GetAll(w http.ResponseWriter, r *http.Request) error {
	var persons []*model.Person
	var err error
	if persons, err = c.PersonRepository.FindAll(); err != nil {
		return err
	}
	if err = c.Render.JSON(w, http.StatusOK, persons); err != nil {
		return err
	}
	return nil
}
