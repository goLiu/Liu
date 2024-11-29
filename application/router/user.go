package router

import (
	"community-governance/application/handlers"
	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes 注册用户路由
func RegisterUserRoutes(r *gin.Engine) {
	userGroup := r.Group("/api/v1/users")
	{
		userGroup.POST("/login", handlers.Login)
	}

}
