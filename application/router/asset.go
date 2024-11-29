package router

import (
	"community-governance/application/handlers"
	"community-governance/application/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterAssetRoutes(r *gin.Engine) {
	assetGroup := r.Group("/api/v1/asset")

	{
		assetGroup.GET("/query/:id", handlers.GetAssetDetail)  // 获取资产信息详细信息
		assetGroup.GET("/query/all", handlers.GetAssetAllPage) // 获取所有资产信息
		assetGroup.GET("/delete/:id", handlers.DeleteAsset)    // 删除资产
		assetGroup.POST("/query/conditions", handlers.GetAssetByConditions)

		recordsGroup := assetGroup.Group("/records") // 修改分组路径
		recordsGroup.Use(middleware.AuthMiddleware())
		{
			recordsGroup.POST("/add", handlers.AddAssetRequestRecord) // 添加记录
			recordsGroup.GET("/query/all", handlers.GetAssetRequestAllPage)
			recordsGroup.GET("/delete/:id", handlers.DeleteAssetRequest)
			recordsGroup.GET("/audit/:id", handlers.AuditRequestRecord)
			recordsGroup.POST("/query/conditions", handlers.GetAssetRequestByConditions)
			recordsGroup.GET("/query/person", handlers.GetAssetRequestByPerson)
		}
	}
}
