package controller

import (
	"fileserver/internal/model"
	"fileserver/internal/repository"
	"github.com/unrolled/render"
	"net/http"
)

type PersonController struct {
	BaseController
	personRepository repository.PersonRepository
}

func (c *PersonController) GetAll(w http.ResponseWriter, r *http.Request) error {
	var persons []*model.Person
	var err error
	if persons, err = c.personRepository.FindAll(); err != nil {
		return err
	}
	if err = c.JsonResponse(w, http.StatusOK, persons); err != nil {
		return err
	}
	return nil
}

func NewPersonController(render *render.Render, personRepository repository.PersonRepository) *PersonController {
	return &PersonController{
		BaseController: BaseController{
			render: render,
		},
		personRepository: personRepository,
	}
}
