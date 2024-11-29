package handlers

import (
	"community-governance/application/models"
	"community-governance/application/utils"
	dbMod "community-governance/db/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

func AddVoteRule(c *gin.Context) {
	var voteRuleReq models.CreateVoteRule
	err := c.ShouldBind(&voteRuleReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "传入的参数不合法:" + err.Error()})
		return
	}

	voteRuleId := uuid.New().String()
	voteRule := dbMod.VoteRule{
		RuleID:      voteRuleId,
		RuleType:    voteRuleReq.RuleType,
		RuleValue:   voteRuleReq.RuleValue,
		Description: voteRuleReq.Description,
		CreateDate:  utils.GetNowTimeString(),
		RuleName:    voteRuleReq.RuleName,
	}
	//添加投票
	err = dbMod.CreateVoteRule(&voteRule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加投票规则失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "添加投票规则成功"})
}

func GetVoteRulesByID(c *gin.Context) {
	//获取路径id值
	id := c.Param("id")
	voteRule, err := dbMod.GetVoteRuleById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取投票规则失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": voteRule})
}

func UpdateVoteRule(c *gin.Context) {
	//获取路径id值
	id := c.Param("id")
	var voteRuleReq dbMod.VoteRule
	err := c.ShouldBind(&voteRuleReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "传入的参数不合法:" + err.Error()})
		return
	}
	err = dbMod.UpdateVoteRule(id, &voteRuleReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新投票规则失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "更新投票规则成功"})

}
func DeleteVoteRule(c *gin.Context) {
	//获取路径id值
	id := c.Param("id")
	err := dbMod.DeleteVoteRule(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除投票规则失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "删除投票规则成功"})
}

func GetVoteRulesByConditions(c *gin.Context) {
	//获取page和pageSize
	page, pageSize := c.Query("page"), c.Query("pageSize")
	iPage, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "传入的page参数不合法:" + err.Error()})
		return
	}
	iPageSize, err := strconv.Atoi(pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "传入的page size参数不合法:" + err.Error()})
		return
	}
	var conditions map[string]interface{}
	err = c.ShouldBind(&conditions)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "传入的参数不合法:" + err.Error()})
		return
	}
	voteRules, err := dbMod.GetVoteRulesAllWithConditions(conditions, iPage, iPageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取投票规则失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": voteRules, "total": len(voteRules)})

}

func GetVoteRulesAllPage(c *gin.Context) {
	//获取page和pageSize
	page, pageSize := c.Query("page"), c.Query("pageSize")
	if page == "" || pageSize == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "传入的page和page size参数不合法"})
		return
	}
	iPage, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "传入的page参数不合法:" + err.Error()})
		return
	}
	iPageSize, err := strconv.Atoi(pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "传入的page size参数不合法:" + err.Error()})
		return
	}
	voteRules, err := dbMod.GetVoteRulesAllWithPagination(iPage, iPageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取投票规则失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": voteRules, "total": len(voteRules)})
}
func GetVoteNames(c *gin.Context) {
	voteNames, err := dbMod.GetVoteRuleNames()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取投票规则失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": voteNames})
}
