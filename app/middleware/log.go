package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

var (
	Logger = logrus.New()
)

func init() {
	logPath := "logs"
	linkName := "logs/log.log"
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		ForceColors:     true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	Logger.Out = os.Stderr
	Logger.SetLevel(logrus.DebugLevel)
	logWriter, _ := rotatelogs.New(
		logPath+"/%Y-%m-%d.log",
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
		rotatelogs.WithLinkName(linkName),
	)
	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}
	Hook := lfshook.NewHook(writeMap, &logrus.TextFormatter{
		FullTimestamp:   true,
		ForceColors:     true,
		TimestampFormat: "2006-01-02 15:04:05", // 格式日志时间
	})
	Logger.SetReportCaller(true)
	Logger.AddHook(Hook)
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//startTime := time.Now()
		uid := strings.ReplaceAll(uuid.New().String(), "-", "")
		c.Set("trace_id", uid)
		c.Next() // 调用该请求的剩余处理程序
		//durationTime := time.Since(startTime).Milliseconds()
		//hostName, err := os.Hostname()
		//if err != nil {
		//    hostName = "Unknown"
		//}
		statusCode := c.Writer.Status()
		//clientIP := c.ClientIP()
		//userAgent := c.Request.UserAgent()
		dataSize := c.Writer.Size()
		if dataSize < 0 {
			dataSize = 0
		}
		method := c.Request.Method
		url := c.Request.RequestURI
		Log := Logger.WithFields(logrus.Fields{
			//"HostName": hostName,
			"url": url,
			"m":   method,
			//"status":   statusCode,
			//"ip":       clientIP,
			"trace_id": uid,
			//"duration": fmt.Sprintf("%d ms", durationTime),
			//"DataSize": dataSize,
			//"UserAgent": userAgent,
		})
		if len(c.Errors) > 0 { // 矿建内部错误
			Log.Error(c.Errors.ByType(gin.ErrorTypePrivate))
		}
		if statusCode >= 500 {
			Log.Error()
		} else if statusCode >= 400 {
			Log.Warn()
		} else {
			Log.Debug()
		}
	}
}
