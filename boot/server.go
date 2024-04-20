package boot

import (
	"github.com/gin-gonic/gin"
	g "main/app/global"
	"main/app/router"
	"net/http"
)

// ServerSetup 函数设置服务器
func ServerSetup() {
	// 从全局配置中获取服务器的配置信息
	config := g.Config.Server

	// 设置gin的模式（debug或release）
	gin.SetMode(config.Mode)
	// 初始化路由
	routers := router.InitRouter()
	// 创建一个HTTP服务器
	server := &http.Server{
		Addr:              config.GetAddr(),         // 服务器地址
		Handler:           routers,                  // 处理器，这里是我们初始化的路由
		TLSConfig:         nil,                      // TLS配置，这里为空，表示不使用TLS
		ReadTimeout:       config.GetReadTimeout(),  // 读取超时时间
		ReadHeaderTimeout: 0,                        // 读取头部超时时间，这里为0，表示没有设置
		WriteTimeout:      config.GetWriteTimeout(), // 写入超时时间
		IdleTimeout:       0,                        // 空闲超时时间，这里为0，表示没有设置
		MaxHeaderBytes:    1 << 20,                  // 最大头部字节数，这里设置为16MB
	}

	// 打印服务器运行的信息
	g.Logger.Infof("server running on %s ...", config.GetAddr())
	// 启动服务器，并在出错时打印错误信息
	g.Logger.Errorf(server.ListenAndServe().Error())
}
