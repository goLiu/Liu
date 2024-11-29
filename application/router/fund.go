package router

import (
	"community-governance/application/handlers"
	"community-governance/application/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterFundRoutes(r *gin.Engine) {
	noticeGroup := r.Group("/api/v1/fund")
	noticeGroup.Use(middleware.AuthMiddleware())
	{
		noticeGroup.POST("/add", handlers.AddFund)             // 创建新财务款项
		noticeGroup.GET("/query/:id", handlers.GetFundDetail)  // 获取财务款项详细信息
		noticeGroup.GET("/query/all", handlers.GetFundAllPage) // 获取所有财务款项信息
		noticeGroup.POST("/query/conditions", handlers.GetFundByConditions)
		noticeGroup.POST("/update/:id", handlers.UpdateFund)        // 更新财务款项信息
		noticeGroup.GET("/delete/:id", handlers.DeleteFund)         // 删除财务款项告信息
		noticeGroup.POST("/add/record/:id", handlers.AddFundRecord) //添加记录
	}

}
