package main

import (
	"github.com/lynxux/goWebFramework/framework/gin"
)

func UserLoginController(c *gin.Context) {
	foo, _ := c.DefaultQueryString("foo", "def")
	//time.Sleep(10 * time.Second)
	c.ISetOkStatus().IJson("ok, UserLoginController" + foo)
}
