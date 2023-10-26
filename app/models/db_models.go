package models

import (
	"github.com/sirupsen/logrus"
	"go-chromedp/app/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
	"time"
)

type MyWriter struct {
	log     *logrus.Logger
	traceid *string
}

func (c *MyWriter) Printf(format string, v ...any) {
	if c.traceid != nil {
		c.log.WithField("trace_id", *c.traceid).Logf(c.log.Level, format, v...)
	} else {
		c.log.Logf(c.log.Level, format, v...)
	}
}

func NewMyWriter() *MyWriter {
	return &MyWriter{
		log: middleware.Logger,
	}
}
func NewMyLogger(traceid string) logger.Interface {
	return logger.New(&MyWriter{
		log:     middleware.Logger,
		traceid: &traceid,
	}, logger.Config{
		Colorful:      true,
		SlowThreshold: 500 * time.Millisecond,
		LogLevel:      logger.Info,
	})
}

var dbInstance *gorm.DB
var once sync.Once

func GetDbInstance() *gorm.DB {
	once.Do(func() {
		dbInstance = dbInit()
	})
	return dbInstance
}
func dbInit() *gorm.DB {
	log := logger.New(NewMyWriter(), logger.Config{
		SlowThreshold: 500 * time.Millisecond,
		Colorful:      true,
		LogLevel:      logger.Info,
	})
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
		Logger: log,
	})
	if err != nil {
		panic(err)
	}
	return db
}
