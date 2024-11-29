package router

import (
	"community-governance/application/handlers"
	"github.com/gin-gonic/gin"
)

// RegisterMemberRoutes 注册成员路由
func RegisterMemberRoutes(r *gin.Engine) {
	memberGroup := r.Group("/api/v1/members")
	{
		memberGroup.GET("query/:id", handlers.GetMember)                 // 获取用户信息
		memberGroup.POST("/add", handlers.AddMember)                     // 创建新用户
		memberGroup.GET("/query/all", handlers.GetMemberAllPage)         // 获取所有用户
		memberGroup.POST("/update/:id", handlers.UpdateMember)           // 更新用户信息
		memberGroup.GET("/delete/:id", handlers.DeleteMember)            // 删除用户
		memberGroup.GET("/update/state/:id", handlers.UpdateMemberState) // 更新用户状态
		memberGroup.POST("/query/conditions", handlers.GetMemberByConditions)
	}
}
