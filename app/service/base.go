package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"go-chromedp/app/middleware"
)

func TraceLog(c *context.Context) *logrus.Entry {
	return middleware.Logger.WithFields(logrus.Fields{
		"trace_id": (*c).Value("trace_id").(string),
	})
}
