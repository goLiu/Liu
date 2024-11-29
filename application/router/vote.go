package router

import (
	"community-governance/application/handlers"
	"community-governance/application/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterVoteRoutes(r *gin.Engine) {
	voteGroup := r.Group("/api/v1/votes")
	voteGroup.Use(middleware.AuthMiddleware())
	{
		voteGroup.POST("/add", handlers.AddVoteProject)                   // 创建新投票项目
		voteGroup.GET("/query/:id", handlers.GetVoteDetail)               // 获取投票项目详细信息
		voteGroup.GET("/query/all", handlers.GetVoteAllPage)              // 获取所有投票项目信息
		voteGroup.POST("/update/:id", handlers.UpdateVote)                // 更新投票项目信息
		voteGroup.GET("/delete/:id", handlers.DeleteVote)                 // 删除投票项目信息
		voteGroup.GET("/update/state/:id", handlers.UpdateVoteState)      // 更新投票状态
		voteGroup.POST("/query/conditions", handlers.GetVoteByConditions) //根据条件查询投票基本信息
		//voteGroup.GET("/query/detail/:id", handlers.GetVoteDetail)        //查询投票详细信息
		voteGroup.GET("/query/join/:id", handlers.VoteJoin) //投票参与
		voteGroup.GET("/end/:id", handlers.VoteEnd)         //投票结束
		// 在 voteGroup 中添加voteRuleGroup子路由组
		voteRuleGroup := voteGroup.Group("/rules")
		{
			voteRuleGroup.POST("/add", handlers.AddVoteRule)                           // 创建投票规则
			voteRuleGroup.GET("/query/:id", handlers.GetVoteRulesByID)                 // 获取投票规则
			voteRuleGroup.GET("/query/all", handlers.GetVoteRulesAllPage)              // 获取所有投票规则
			voteRuleGroup.POST("/update/:id", handlers.UpdateVoteRule)                 // 更新投票规则
			voteRuleGroup.GET("/delete/:id", handlers.DeleteVoteRule)                  // 删除投票规则
			voteRuleGroup.POST("/query/conditions", handlers.GetVoteRulesByConditions) //根据条件查询投票规则
			voteRuleGroup.GET("/query/names", handlers.GetVoteNames)                   // 获取所有投票规则的名称
		}
	}

}
