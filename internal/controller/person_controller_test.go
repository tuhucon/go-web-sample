package controller

import (
	"encoding/json"
	"fileserver/internal/model"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/mock"
	"github.com/unrolled/render"
	"net/http"
	"net/http/httptest"
	"testing"
)

type PersonRepositoryMock struct {
	mock.Mock
}

func (p *PersonRepositoryMock) FindAll() ([]*model.Person, error) {
	args := p.Called()
	return args.Get(0).([]*model.Person), args.Error(1)
}

func TestPersonController_GetAll_ReturnEmpty(t *testing.T) {
	personRepositoryMock := new(PersonRepositoryMock)
	personRepositoryMock.On("FindAll").Return([]*model.Person{}, nil)
	personController := NewPersonController(render.New(), personRepositoryMock)

	r := httptest.NewRequest(http.MethodGet, "/persons", nil)
	w := httptest.NewRecorder()
	err := personController.GetAll(w, r)

	personRepositoryMock.AssertExpectations(t)
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(w.Code, http.StatusOK); diff != "" {
		t.Error(diff)
	}
	result := []*model.Person{}
	err = json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(0, len(result)); diff != "" {
		t.Error(diff)
	}
}

func TestPersonController_GetAll_ReturnOne(t *testing.T) {
	dataTest := []*model.Person{{1, "tu hu con", 12}}
	personRepositoryMock := new(PersonRepositoryMock)
	personRepositoryMock.On("FindAll").Return(dataTest, nil)
	personController := NewPersonController(render.New(), personRepositoryMock)

	r := httptest.NewRequest(http.MethodGet, "/persons", nil)
	w := httptest.NewRecorder()
	err := personController.GetAll(w, r)

	personRepositoryMock.AssertExpectations(t)
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(w.Code, http.StatusOK); diff != "" {
		t.Error(diff)
	}
	var result []*model.Person
	err = json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(dataTest, result); diff != "" {
		t.Error(diff)
	}
}
