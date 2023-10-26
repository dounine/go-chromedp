package views

import (
	"github.com/gin-gonic/gin"
	"go-chromedp/app/service"
)

type DumpRouter struct {
	dumpService service.IDumpService
}

func NewDumpRouter() *DumpRouter {
	return &DumpRouter{
		dumpService: service.NewDumpService(),
	}
}

func (*DumpRouter) Init(r *gin.Engine) {
	defaultRouter := r.Group("/dump")
	controller := NewDumpController()
	g := GetReturnHandlerFunc()
	defaultRouter.GET("/info", g(controller.Info))
	//defaultRouter.GET("/:p8id/:id", P8IDIDValidMiddleware(), g(controller.Info))
	//defaultRouter.PATCH("/:p8id/:id/disable", P8IDIDValidMiddleware(), g(controller.Disable))
	//defaultRouter.PATCH("/:p8id/:id/enable", P8IDIDValidMiddleware(), g(controller.Enable))
	//defaultRouter.POST("/:p8id", g(controller.Create))
	//defaultRouter.PATCH("/:p8id/:id", P8IDIDValidMiddleware(), g(controller.Update))
}
