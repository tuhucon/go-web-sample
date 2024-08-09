package controller

import (
	"github.com/unrolled/render"
	"net/http"
)

type BaseController struct {
	Render *render.Render
}

func (c *BaseController) JsonResponse(w http.ResponseWriter, statusCode int, data any) error {
	return c.Render.JSON(w, statusCode, data)
}
