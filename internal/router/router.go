package router

import (
	"time"

	"openapi/internal/handler"
	"openapi/internal/logger"
	"openapi/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter 设置 API 路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 添加全局中间件
	r.Use(gin.Recovery())
	r.Use(corsMiddleware())
	r.Use(requestLogger())

	// 健康检查路由
	r.GET("/healthcheck", handler.HealthCheckHandler)

	// Swagger 文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API 路由组
	api := r.Group("/api/v1")
	{
		// 认证相关路由组
		auth := api.Group("/auth")
		{
			auth.POST("/login", handler.LoginHandler)
			auth.POST("/register", handler.RegisterHandler)
			auth.POST("/logout", middleware.JWTAuth(), handler.LogoutHandler)
		}

		// 需要认证的路由
		protected := api.Group("")
		protected.Use(middleware.JWTAuth())
		{
			// FC服务管理路由组
			services := protected.Group("/services")
			{
				services.GET("", middleware.ValidateAndGetHeaders(middleware.CommonHeaders...), handler.ListServiceHandler)
				services.POST("/versions", middleware.ValidateAndGetHeaders(middleware.CommonHeaders...), handler.ListServiceVersionHandler)
				services.POST("/functions", middleware.ValidateAndGetHeaders(middleware.CommonHeaders...), handler.ListFcHandler)
				services.POST("/aliases", middleware.ValidateAndGetHeaders(middleware.CommonHeaders...), handler.ListAliasHandler)
				services.PUT("/aliases", middleware.ValidateAndGetHeaders(middleware.CommonHeaders...), handler.UpdateAliasHandler)
				services.POST("/publish", middleware.ValidateAndGetHeaders(middleware.CommonHeaders...), handler.PublicServiceHandler)
			}

			// 系统管理路由组
			system := protected.Group("/system")
			{
				system.GET("/zones", handler.GetZoneInfoHandler)
			}

			//Cloudflare接口管理路由组
			cloudflare := protected.Group("/cloudflare")
			{
				cloudflare.GET("/pages/info", middleware.ValidateAndGetHeaders(middleware.CommonHeaders...), handler.GetPagesProjectHandlerInfo)
				cloudflare.GET("/pages/projects", middleware.ValidateAndGetHeaders(middleware.CommonHeaders...), handler.GetPagesProjectHandler)
				cloudflare.GET("/kv/namespaces", middleware.ValidateAndGetHeaders(middleware.CommonHeaders...), handler.GetKVNamespacesHandler)
				cloudflare.POST("/kv/namespaces/keys", middleware.ValidateAndGetHeaders(middleware.CommonHeaders...), handler.GetKVKeysHandler)
				cloudflare.POST("/kv/namespaces/keys/values", middleware.ValidateAndGetHeaders(middleware.CommonHeaders...), handler.GetKVKeyValuesHandler)
				cloudflare.PUT("/kv/namespaces/keys/values", middleware.ValidateAndGetHeaders(middleware.CommonHeaders...), handler.UpdateKVKeyValuesHandler)
				cloudflare.POST("/bucketinfo", middleware.ValidateAndGetHeaders(middleware.CommonHeaders...), handler.GetBucketHandler)
				cloudflare.DELETE("/bucketinfo", middleware.ValidateAndGetHeaders(middleware.CommonHeaders...), handler.DeleteDirectoryHandler)
				cloudflare.POST("/bucketinfo/copy", middleware.ValidateAndGetHeaders(middleware.CommonHeaders...), handler.CopyDirectoryHandler)

				// Cloudflare账号管理路由组
				accounts := cloudflare.Group("/accounts")
				{
					accounts.GET("", handler.ListCloudflareAccountsHandler)
					accounts.POST("", handler.CreateCloudflareAccountHandler)
					accounts.PUT("/:id", handler.UpdateCloudflareAccountHandler)
					accounts.DELETE("/:id", handler.DeleteCloudflareAccountHandler)
				}
			}
		}
	}

	return r
}

// corsMiddleware 处理跨域请求
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Env, Country-Code")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// requestLogger 记录请求日志
func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		endTime := time.Now()
		latency := endTime.Sub(startTime)

		logger.Info("Request: %s %s %d %v",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			latency,
		)
	}
}
