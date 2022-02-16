package middleware

import (
	"github.com/lynxux/goWebFramework/framework/gin"
	"log"
	"time"
)

// 记录中间件的运行时间
func Cost() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		end := time.Now()
		cost := end.Sub(start)
		log.Printf("api uri: %v, cost: %v", c.Request.RequestURI, cost.Seconds())
	}
}
