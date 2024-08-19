package controller

import (
	"encoding/json"
	"fileserver/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/mock"
	"github.com/unrolled/render"
)

type PersonRepositoryMock struct {
	mock.Mock
}

func (p *PersonRepositoryMock) FindAll() ([]*model.Person, error) {
	args := p.Called()
	return args.Get(0).([]*model.Person), args.Error(1)
}

func TestPersonController_GetAll(t *testing.T) {
	table := []struct {
		name string
		data []*model.Person
	}{
		{"no data", []*model.Person{}},
		{"one person", []*model.Person{{1, "tu hu con", 20}}},
		{"two persons", []*model.Person{{1, "tu hu con", 20}, {2, "chich choe", 30}}},
	}
	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			personRepositoryMock := new(PersonRepositoryMock)
			personRepositoryMock.On("FindAll").Return(tt.data, nil)
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
			if diff := cmp.Diff(tt.data, result); diff != "" {
				t.Error(diff)
			}
		})
	}
}
