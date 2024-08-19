package controller

import (
	"net/http"

	"github.com/unrolled/render"
)

// BaseController struct provides JSONResponse method
// All your controller have to embed BaseController and use JSONResponse method to render Json
type BaseController struct {
	render *render.Render
}

// JSONResponse method return your data as json
// Behind the scene, JSONResponse method use render.Render.JSON
func (c *BaseController) JSONResponse(w http.ResponseWriter, statusCode int, data any) error {
	return c.render.JSON(w, statusCode, data)
}
