package framework

import (
	"context"
	"fmt"
	"log"
	"time"
)

func TimeoutHandler(fun ControllerHandler, d time.Duration) ControllerHandler {
	// 使用函数回调
	return func(c *Context) error {
		finish := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)
		// 执行业务逻辑前预操作：初始化超时context
		durationCtx, cancel := context.WithTimeout(c.BaseContext(), d)
		defer cancel()

		c.request.WithContext(durationCtx)

		// golang 中每个Goroutine创建时都需要使用defer和recover对该线程可能发生的panic做处理，否则任意一个线程的panic都可能导致进程的panic
		go func() {
			defer func() { // panic做处理
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			// 业务处理
			fun(c)
			// 通知结束
			finish <- struct{}{}
		}()
		//监听事件
		select {
		case p := <-panicChan:
			c.WriterMux().Lock()         // 当异常事件发生时，可能需要对responseWriter中写入数据
			defer c.WriterMux().Unlock() // 为了防止其他的routine对responseWriter的访问，需要使用锁
			log.Println(p)
			c.SetStatus(500).Json("panic")
		case <-finish:
			fmt.Println("finish")
		case <-durationCtx.Done():
			c.WriterMux().Lock() // 当超时事件发生时，也同理
			defer c.WriterMux().Unlock()
			c.SetStatus(500).Json("time out")
			c.SetHasTimeout() // 当超时事件发生之后，已经向responseWriter中写入数据，为了方式其他Goroutine再次写入数据，需要设置超时标志位
		}
		return nil
	}
}
