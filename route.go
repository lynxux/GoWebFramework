package main

import (
	"github.com/lynxux/goWebFramework/framework/gin"
	"github.com/lynxux/goWebFramework/framework/middleware"
)

// 注册路由规则
func registerRouter(core *gin.Engine) {
	// 静态路由+HTTP方法匹配
	core.GET("/user/login", middleware.Test3(), UserLoginController) // 单个路由增加中间件
	// 批量通用前缀
	subjectApi := core.Group("/subject")
	{
		// 动态路由
		subjectApi.DELETE("/:id", SubjectDelController)
		subjectApi.PUT("/:id", SubjectUpdateController)
		// 单个路由增加中间件
		subjectApi.GET("/:id", middleware.Test3(), SubjectGetController)
		subjectApi.GET("/list/all", SubjectListController)
	}
}
