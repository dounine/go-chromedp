package views

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go-chromedp/app/middleware"
	"strconv"
)

type (
	CRUD interface {
		List(c *gin.Context)
		Create(c *gin.Context)
		Update(c *gin.Context)
		Delete(c *gin.Context)
		Info(c *gin.Context)
	}

	ResponseData struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data any    `json:"data"`
	}
	Paging struct {
		Offset      int     `form:"offset" binding:"numeric,gte=0"`
		Limit       int     `form:"limit" binding:"numeric,gt=0,lte=20"`
		FilterKey   *string `form:"filter_key"`
		FilterValue *string `form:"filter_value" binding:"required_with=FilterKey"`
	}
	P8IDAndID struct {
		P8ID int    `json:"p8id" uri:"p8id" binding:"required"`
		ID   string `json:"id" uri:"id" binding:"required"`
	}
)

func PagingValidMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		paging := Paging{
			Limit: 20,
		}
		if err := c.ShouldBindQuery(&paging); err != nil {
			c.AbortWithStatusJSON(200, gin.H{
				"msg": err.Error() + ":" + "参数校验失败",
			})
			return
		}
		c.Next()
	}
}
func P8IDValidMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		p8id := c.Param("p8id")
		_, err := strconv.Atoi(p8id)
		if err != nil {
			c.AbortWithStatusJSON(200, gin.H{
				"msg": err.Error() + ":" + "参数校验失败",
			})
			return
		}
		c.Next()
	}
}
func IntIDValidMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		intId := c.Param("id")
		_, err := strconv.Atoi(intId)
		if err != nil {
			c.AbortWithStatusJSON(200, gin.H{
				"msg": err.Error() + ":" + "参数校验失败",
			})
			return
		}
		c.Next()
	}
}
func P8IDIDValidMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		p8idid := P8IDAndID{}
		if err := c.ShouldBindUri(&p8idid); err != nil {
			c.AbortWithStatusJSON(200, gin.H{
				"msg": err.Error() + ":" + "参数校验失败",
			})
			return
		}
		c.Next()
	}
}

func TraceContext(c *gin.Context) context.Context {
	ctx := c.Request.Context()
	cc := context.WithValue(ctx, "trace_id", c.MustGet("trace_id"))
	return cc
}
func TraceLog(ctx *context.Context) *logrus.Entry {
	return middleware.Logger.WithFields(logrus.Fields{
		"trace_id": (*ctx).Value("trace_id").(string),
	})
}

func fail(g *gin.Context, err error) {
	g.JSON(200, gin.H{
		"msg": err.Error(),
	})
}
func ok(g *gin.Context, data any) {
	if data == nil {
		g.JSON(200, gin.H{
			"code": 1,
		})
	} else {
		g.JSON(200, gin.H{
			"code": 1,
			"data": data,
		})
	}
}

type ReturnHandlerFunc func(c *gin.Context) (any, error)

func GetReturnHandlerFunc() func(ReturnHandlerFunc) gin.HandlerFunc {
	return func(handlerFunc ReturnHandlerFunc) gin.HandlerFunc {
		return handlerFunc.GinHandlerFunc()
	}
}
func (gh ReturnHandlerFunc) GinHandlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := gh(c)
		if err != nil {
			fail(c, err)
			return
		}
		ok(c, data)
	}
}
