package router

import (
	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	RegisterFundRoutes(r)
	RegisterNoticeRoutes(r)
	RegisterUserRoutes(r)
	RegisterAssetRoutes(r)
	RegisterMemberRoutes(r)
	RegisterVoteRoutes(r)
	RegisterFacilityRoutes(r)
	return r
}
