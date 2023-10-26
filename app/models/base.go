package models

import (
	"context"
	"github.com/sirupsen/logrus"
	"go-chromedp/app/middleware"
)

type (
	TupleData interface {
		any | []any
	}
	ResponseData[T TupleData] struct {
		Errors []struct {
			Status string `json:"status"`
			Code   string `json:"code"`
			Title  string `json:"title"`
			Detail string `json:"detail"`
		}
		Data T     `json:"data"`
		Meta *Meta `json:"meta"`
	}
	Paging struct {
		Total int `json:"total"`
		Limit int `json:"limit"`
	}
	Meta struct {
		Paging Paging `json:"paging"`
	}
)

func errCodes(codes ...int) map[int]bool {
	m := map[int]bool{}
	for _, code := range codes {
		m[code] = true
	}
	return m
}

func TraceLog(ctx *context.Context) *logrus.Entry {
	return middleware.Logger.WithFields(logrus.Fields{
		"trace_id": (*ctx).Value("trace_id").(string),
	})
}
