package controller

import (
	"net/http"
	"time"

	"github.com/unrolled/render"
)

type HelloController struct {
	BaseController
}

func (h *HelloController) Hello(w http.ResponseWriter, r *http.Request) error {
	return h.JSONResponse(w, http.StatusOK, "Hello World")
}

func (h *HelloController) Time(w http.ResponseWriter, r *http.Request) error {
	now := time.Now()
	return h.JSONResponse(w, http.StatusOK, now.Format(time.DateTime))
}

func NewHelloController(render *render.Render) *HelloController {
	return &HelloController{
		BaseController{
			render: render,
		},
	}
}
