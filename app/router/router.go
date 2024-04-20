package router

import (
	"github.com/gin-gonic/gin"
	g "main/app/global"
	"main/app/internal/middleware"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// 使用中间件，包括Zap日志记录器、Zap恢复和按规则的跨域资源共享
	r.Use(middleware.ZapLogger(g.Logger), middleware.ZapRecovery(g.Logger, true))
	r.Use(middleware.CorsByRules())

	// 创建一个新的路由组
	routerGroup := new(Group)

	// 创建一个公共的路由组
	PublicGroup := r.Group("/api")
	{
		routerGroup.InitUserSignRouter(PublicGroup)
	}

	// 创建一个私有的路由组，并使用JWT认证中间件
	PrivateGroup := r.Group("/api")
	PrivateGroup.Use(middleware.JWTAuthMiddleware())
	{
		routerGroup.InitRecipeRouter(PrivateGroup)
		routerGroup.InitRestaurantRouter(PrivateGroup)
		routerGroup.InitUserInfoRouter(PrivateGroup)
	}

	// 记录初始化路由成功的信息
	g.Logger.Infof("initialize routers successfully")
	return r
}
