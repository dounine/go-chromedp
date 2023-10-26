package store

import (
	"context"
	"go-chromedp/app/models"
	"gorm.io/gorm"
)

func TraceDB(ctx *context.Context) *gorm.DB {
	return models.GetDbInstance().Session(&gorm.Session{
		Context: *ctx,
		Logger:  models.NewMyLogger((*ctx).Value("trace_id").(string)),
	})
}

func InitSchemas() {
	NewCodeStore().InitSchema()
	NewCodeUseStore().InitSchema()
}
