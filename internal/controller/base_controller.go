package controller

import (
	"github.com/unrolled/render"
	"net/http"
)

type BaseController struct {
	render *render.Render
}

func (c *BaseController) JsonResponse(w http.ResponseWriter, statusCode int, data any) error {
	return c.render.JSON(w, statusCode, data)
}
