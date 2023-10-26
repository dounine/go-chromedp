package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"golang.org/x/time/rate"
	"time"
)

func CreateIpLimiter() *rate.Limiter {
	return rate.NewLimiter(rate.Every(time.Second*1), 20)
}
func IPLimiter() gin.HandlerFunc {
	ipCaches := expirable.NewLRU[string, *rate.Limiter](1000, nil, time.Hour*24)
	tooManyRequest := gin.H{"msg": "请求过于频繁"}
	return func(c *gin.Context) {
		ip := c.ClientIP()
		var limiter *rate.Limiter
		var ok bool
		if limiter, ok = ipCaches.Get(ip); !ok {
			limiter = CreateIpLimiter()
			ipCaches.Add(ip, limiter)
		}
		if limiter.Allow() {
			c.Next()
		} else {
			Logger.Error("too many request ", ip)
			c.AbortWithStatusJSON(200, tooManyRequest)
		}
	}
}
