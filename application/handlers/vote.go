package handlers

import (
	"community-governance/application/models"
	"community-governance/application/utils"
	dbMod "community-governance/db/models"
	"community-governance/fabric"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"strings"
)

const (
	VoteStateActive   = "active"
	VoteStateInactive = "inactive"
	VoteStateEnd      = "end"
)

func AddVoteProject(c *gin.Context) {
	var voteReq models.CreateVote
	err := c.ShouldBind(&voteReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "传入的参数不合法:" + err.Error()})
		return
	}
	//获取规则id
	ruleId, err := dbMod.GetRuleIdByName(voteReq.RuleName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取规则id失败:" + err.Error()})
		return
	}
	voteId := uuid.New().String()
	vote := dbMod.Vote{
		VoteID:      voteId,
		Name:        voteReq.Name,
		RuleID:      ruleId,
		StartTime:   voteReq.StartTime,
		Manager:     voteReq.Manager,
		Description: voteReq.Description,
		Status:      VoteStateActive,
	}
	//添加投票
	err = dbMod.CreateVote(&vote)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加投票失败:" + err.Error()})
		return
	}
	nowTime := utils.GetNowTimeString()
	optionsStr := make([]string, len(voteReq.Options))
	for _, opVal := range voteReq.Options {
		option := dbMod.VoteOption{
			VoteID:      voteId,
			OptionID:    uuid.New().String(),
			OptionValue: opVal.OptionValue,
			Status:      VoteStateActive,
			Description: opVal.Description,
			CreateDate:  nowTime,
		}
		optionsStr = append(optionsStr, opVal.OptionValue)
		//创建option
		err := dbMod.CreateVoteOption(&option)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "添加option失败:" + err.Error()})
			return
		}
	}
	//计算hash值
	hash, err := utils.ComputeHash(vote)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "计算hash值失败:" + err.Error()})
		return
	}
	//获取ruleType
	rule, err := dbMod.GetVoteRuleById(ruleId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取ruleType失败:" + err.Error()})
		return
	}
	//调用合约
	err = fabric.CreatVote(voteId, hash, rule.RuleType, rule.RuleValue, strings.Join(optionsStr, ","))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "调用合约失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "添加投票成功"})
}

// UpdateVote 更新投票项目信息
func UpdateVote(c *gin.Context) {
	var voteReq dbMod.Vote
	//获取路径id值
	id := c.Param("id")
	err := c.ShouldBind(&voteReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "传入的参数不合法:" + err.Error()})
		return
	}
	err = dbMod.UpdateVote(id, &voteReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新投票失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "更新投票成功"})
}

func GetVoteDetail(c *gin.Context) {
	type Detail struct {
		Base dbMod.Vote `json:"base"`
		fabric.ChainVoteDetail
		Options []dbMod.VoteOption `json:"options"`
	}
	var detail Detail
	//获取路径id值
	id := c.Param("id")
	base, err := dbMod.GetVoteByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取投票失败:" + err.Error()})
		return
	}
	detail.Base = *base
	options, err := dbMod.GetVoteOptionByVoteId(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取投票选项失败:" + err.Error()})
		return
	}
	detail.Options = options
	chainDetail, err := fabric.GetVoteRecordDetail(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取投票详情失败:" + err.Error()})
		return
	}
	detail.ChainVoteDetail = chainDetail
	c.JSON(http.StatusOK, gin.H{"data": detail})
}

func GetVoteAllPage(c *gin.Context) {
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
	votes, err := dbMod.GetVoteAllWithPagination(iPage, iPageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取投票失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": votes, "total": len(votes)})
}

func DeleteVote(c *gin.Context) {
	//获取路径id值
	id := c.Param("id")
	err := dbMod.DeleteVote(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除投票失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "删除投票成功"})

}

func UpdateVoteState(c *gin.Context) {
	//获取路径id值
	id := c.Param("id")
	state := c.Param("state")
	updateMap := map[string]interface{}{
		"status": state,
	}
	err := dbMod.UpdateVote(id, updateMap)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新投票状态失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "更新投票状态成功"})
}

func GetVoteByConditions(c *gin.Context) {
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
	votes, err := dbMod.GetVoteAllWithConditions(conditions, iPage, iPageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取投票失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": votes, "total": len(votes)})
}

// VoteJoin 参与投票
func VoteJoin(c *gin.Context) {
	//获取路径id值
	id := c.Param("id")
	//获取option与voter
	option := c.Query("option")
	if option == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "传入的参数不合法"})
		return
	}
	//获取userId
	userId := c.MustGet("userId").(string)
	result, err := fabric.VoteJoin(id, userId, option)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "参与投票失败:" + err.Error()})
		return
	}
	//若不为空，则说明投票结束
	if result != "" {
		//更新投票状态
		updateMap := map[string]interface{}{
			"status": VoteStateEnd,
			"result": result,
		}
		err := dbMod.UpdateVote(id, updateMap)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新投票状态失败:" + err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"data": "参与投票成功"})
}
func VoteEnd(c *gin.Context) {
	//获取路径id值
	id := c.Param("id")
	result, err := fabric.EndVote(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "结束投票失败:" + err.Error()})
		return
	}
	//更新投票状态
	updateMap := map[string]interface{}{
		"status": VoteStateEnd,
		"result": strings.Join(result, ","),
	}
	err = dbMod.UpdateVote(id, updateMap)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新投票状态失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "结束投票成功"})
}
