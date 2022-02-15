package middleware

import (
	"github.com/lynxux/goWebFramework/framework"
	"log"
	"time"
)

// 记录中间件的运行时间
func Cost() framework.ControllerHandler {
	return func(c *framework.Context) error {
		start := time.Now()

		c.Next()

		end := time.Now()
		cost := end.Sub(start)
		log.Printf("api uri: %v, cost: %v", c.GetRequest().RequestURI, cost.Seconds())
		return nil
	}
}
