package handlers

import (
	"community-governance/application/models"
	"community-governance/application/utils"
	dbMod "community-governance/db/models"
	"community-governance/fabric"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

const (
	foundStateActive    = "active"
	foundStateInactive  = "inactive"
	foundStateClosed    = "closed"
	foundStatePending   = "pending"
	foundStateCancelled = "cancelled"
)

func AddFund(c *gin.Context) {
	var fundReq models.CreateFund
	if err := c.ShouldBindJSON(&fundReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fundId := uuid.New().String()
	fund := dbMod.Fund{
		FundID:            fundId,
		Description:       fundReq.Description,
		Status:            foundStateActive,
		Name:              fundReq.Name,
		Source:            fundReq.Source,
		SourceDescription: fundReq.SourceDescription,
		TotalAmount:       fundReq.TotalAmount,
		CurrentBalance:    fundReq.TotalAmount,
		Manager:           fundReq.Manager,
		EstablishDate:     fundReq.EstablishDate,
	}
	if err := dbMod.CreateFunds(&fund); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//计算hash值
	hash, err := utils.ComputeHash(fund)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "计算公告hash失败:" + err.Error()})
		return
	}
	err = fabric.CreateFinancial(fundId, hash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新公告失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "创建财务款项成功"})
}

func GetFundDetail(c *gin.Context) {
	type FundDetail struct {
		Base dbMod.Fund
		fabric.ChainFundDetail
	}
	var detail FundDetail
	//获取路径id值
	id := c.Param("id")
	base, err := dbMod.GetFundByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取投票失败:" + err.Error()})
		return
	}
	detail.Base = *base
	chainDetail, err := fabric.GetFinancialDetail(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取投票详情失败:" + err.Error()})
		return
	}
	detail.ChainFundDetail = chainDetail
	c.JSON(http.StatusOK, gin.H{"data": detail})
}

func GetFundAllPage(c *gin.Context) {
	//获取page和pageSize
	page, pageSize := c.Query("page"), c.Query("pageSize")
	if page == "" || pageSize == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "传入的page或pageSize参数不合法"})
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
	funds, err := dbMod.GetAllFundsWithPagination(iPage, iPageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取公告失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": funds, "total": len(funds)})

}
func UpdateFund(c *gin.Context) {
	//获取路径id值
	id := c.Param("id")
	var fundReq dbMod.Fund
	if err := c.ShouldBindJSON(&fundReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := dbMod.UpdateFund(id, fundReq); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新公告失败:" + err.Error()})
		return
	}
	//计算hash值
	hash, err := utils.ComputeHash(fundReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "计算公告hash失败:" + err.Error()})
		return
	}
	err = fabric.UpdateFinancial(id, hash, fundReq.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新公告失败:" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "更新公告成功"})
}
func DeleteFund(c *gin.Context) {
	//获取路径id值
	id := c.Param("id")
	err := dbMod.DeleteFund(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除公告失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "删除公告成功"})
}
func AddFundRecord(c *gin.Context) {
	//获取路径id值
	id := c.Param("id")
	var recordReq models.FinancialRecord
	err := c.ShouldBind(&recordReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "传入的参数不合法:" + err.Error()})
		return
	}
	//获取fund信息
	fund, err := dbMod.GetFundByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取公告失败:" + err.Error()})
		return
	}
	currentBalance, err := strconv.ParseFloat(fund.CurrentBalance, 64)
	if err != nil {
		fmt.Println("转换current balance错误:", err)
		return
	}
	recordAmount, err := strconv.ParseFloat(recordReq.Amount, 64)
	if err != nil {
		fmt.Println("转换record amount错误:", err)
	}
	if recordReq.Type == "1" { //支出
		if currentBalance < recordAmount {
			c.JSON(http.StatusBadRequest, gin.H{"error": "余额不足"})
			return
		}
		recordAmount = -recordAmount
	}
	//计算fund剩余余额
	nowBalance := currentBalance + recordAmount
	nowBalanceStr := strconv.FormatFloat(nowBalance, 'f', -1, 64)
	if err := dbMod.UpdateFund(id, *fund); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新公告失败:" + err.Error()})
		return
	}
	//更新记录
	updateMap := map[string]interface{}{
		"current_balance": nowBalanceStr,
	}
	err = dbMod.UpdateFund(fund.FundID, updateMap)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新公告失败:" + err.Error()})
		return
	}
	userId := c.MustGet("userId").(string)
	err = fabric.AddFinancialRecord(id, recordReq.Type, recordReq.Source, recordReq.Explain, userId, recordReq.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加记录失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "添加记录成功"})
}

func GetFundByConditions(c *gin.Context) {
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
	funds, err := dbMod.GetAllFundsWithConditions(conditions, iPage, iPageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取公告失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": funds, "total": len(funds)})
}
