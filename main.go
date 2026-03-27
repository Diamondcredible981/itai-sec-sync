package main

import (
	"github.com/gin-gonic/gin"
	"github.com/iMayday-Yee/XinchuangAnalyze/database"
	"github.com/iMayday-Yee/XinchuangAnalyze/middleware"
	"github.com/iMayday-Yee/XinchuangAnalyze/routers"
)

// 主函数
func main() {
	db := database.SetupDatabase()
	router := gin.Default()
	router.Use(middleware.CORS())
	router.Use(middleware.AuditLogger())
	routers.RegisterRouters(router, db)
	router.Run(":8080")
}
