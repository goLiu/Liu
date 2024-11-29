package router

import (
	"community-governance/application/handlers"
	"community-governance/application/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterFacilityRoutes(r *gin.Engine) {
	faclitiesGroup := r.Group("/api/v1/facilities")
	faclitiesGroup.Use(middleware.AuthMiddleware())
	{
		faclitiesGroup.GET("/query/:id", handlers.GetFacilityDetail)               // 获取设备信息详细信息
		faclitiesGroup.POST("/add", handlers.AddFacility)                          // 创建新公共设施
		faclitiesGroup.GET("/query/all", handlers.GetFacilityAllPage)              // 获取所有公共设施
		faclitiesGroup.POST("/update/:id", handlers.UpdateFacility)                // 更新设备信息
		faclitiesGroup.GET("/delete/:id", handlers.DeleteFacility)                 // 删除设备
		faclitiesGroup.POST("/query/conditions", handlers.GetFacilityByConditions) // 根据条件获取公共设施
		faclitiesGroup.GET("/request/:id", handlers.RequestFacility)               // 请求公共设施
		faclitiesGroup.GET("/release/:id", handlers.ReleaseFacility)               // 释放公共设施
	}
}
