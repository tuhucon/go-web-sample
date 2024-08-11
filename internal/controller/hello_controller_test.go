package controller

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/unrolled/render"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var helloController = NewHelloController(render.New())

func TestHelloController_Hello(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/hello", nil)
	w := httptest.NewRecorder()
	err := helloController.Hello(w, r)
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(w.Code, http.StatusOK); diff != "" {
		t.Error(diff)
	}
	if diff := cmp.Diff(w.Body.String(), "\"Hello World\""); diff != "" {
		t.Error(diff)
	}
}

func TestHelloController_Time(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/time", nil)
	w := httptest.NewRecorder()
	err := helloController.Time(w, r)
	datetime := time.Now().Format(time.DateTime)
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(w.Code, http.StatusOK); diff != "" {
		t.Error(diff)
	}
	if diff := cmp.Diff(w.Body.String(), fmt.Sprintf("\"%s\"", datetime)); diff != "" {
		t.Error(diff)
	}
}
