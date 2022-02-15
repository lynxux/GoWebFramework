package main

import (
	"context"
	"fmt"
	"github.com/lynxux/goWebFramework/framework"
	"log"
	"time"
)

func FooControllerHandler(c *framework.Context) error {
	// 在业务逻辑处理前，创建有定时器功能的 context
	durationContext, cancel := context.WithTimeout(c.BaseContext(), time.Duration(1*time.Second))
	defer cancel()

	// 用于通知结束
	finish := make(chan struct{}, 1)
	// 用于通知panic异常
	panicChan := make(chan interface{}, 1)
	// golang 中每个Goroutine创建时都需要使用defer和recover对该线程可能发生的panic做处理，否则任意一个线程的panic都可能导致进程的panic
	go func() {
		defer func() { // panic做处理
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()
		// 业务处理
		time.Sleep(10 * time.Second)
		c.SetOkStatus().Json("ok")
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
	case <-durationContext.Done():
		c.WriterMux().Lock() // 当超时事件发生时，也同理
		defer c.WriterMux().Unlock()
		c.SetStatus(500).Json("timeout")
		c.SetHasTimeout() // 当超时事件发生之后，已经向responseWriter中写入数据，为了方式其他Goroutine再次写入数据，需要设置超时标志位
	}
	return nil
}
