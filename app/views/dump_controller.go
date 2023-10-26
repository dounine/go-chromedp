package views

import "github.com/gin-gonic/gin"

type DumpController struct {
}

func NewDumpController() *DumpController {
	return &DumpController{}
}
func (c *DumpController) Info(*gin.Context) (any, error) {
	return gin.H{
		"version": "1.0.0",
	}, nil
}
