package router

import (
	"community-governance/application/handlers"
	"community-governance/application/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterNoticeRoutes(r *gin.Engine) {
	noticeGroup := r.Group("/api/v1/notice")
	noticeGroup.Use(middleware.AuthMiddleware())
	{
		noticeGroup.POST("/add", handlers.AddNotice)             // 创建新公告
		noticeGroup.GET("/query/:id", handlers.GetNoticeDetail)  // 获取公告详细信息
		noticeGroup.GET("/query/all", handlers.GetNoticeAllPage) // 获取所有公告信息
		noticeGroup.POST("/update/:id", handlers.UpdateNotice)   // 更新公告信息
		noticeGroup.GET("/delete/:id", handlers.DeleteNotice)    // 删除公告信息
		noticeGroup.POST("/query/conditions", handlers.GetNoticeByConditions)
	}

}
