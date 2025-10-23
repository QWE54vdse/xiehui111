package main

import (
	"go-auth/database"
	"go-auth/handlers"
	"go-auth/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 连接数据库
	database.ConnectDB()

	// 自动迁移数据库表
	db := database.GetDB()
	// db.Migrator().DropTable(&models.User{})
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		panic("数据库迁移失败: " + err.Error())
	}

	// 创建Gin路由
	r := gin.Default()

	// 简化CORS配置
	r.Use(cors.Default()) // 使用默认配置，允许所有来源

	// 或者使用更宽松的配置
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true, // 允许所有来源
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 路由定义
	api := r.Group("/api")
	{
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)

		auth := api.Group("/auth")
		auth.Use(handlers.AuthMiddleware())
		{
			auth.GET("/profile", handlers.GetProfile)
		}
	}

	// 启动服务器
	r.Run(":8080")
}
