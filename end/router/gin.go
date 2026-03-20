package router

import "github.com/gin-gonic/gin"

func init() {
	// 禁用 Gin 在启动时打印路由表（[GIN-debug] ... -> ...）
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {}
}
