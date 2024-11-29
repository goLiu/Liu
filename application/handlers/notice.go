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
)

func AddNotice(c *gin.Context) {
	var noticeReq models.CreateNotice
	err := c.ShouldBind(&noticeReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "传入的参数不合法:" + err.Error()})
		return
	}
	noticeId := uuid.New().String()
	userId := c.MustGet("userId").(string)
	notice := dbMod.Notice{
		NoticeID:    noticeId,
		Title:       noticeReq.Title,
		Content:     noticeReq.Content,
		Type:        noticeReq.Type,
		Author:      userId,
		PublishTime: utils.GetNowTimeString(),
		Version:     1,
	}
	err = dbMod.CreateNotice(notice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加公告失败:" + err.Error()})
		return
	}
	//计算公告hash
	hash, err := utils.ComputeHash(notice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "计算公告hash失败:" + err.Error()})
		return
	}
	err = fabric.CreateNotice(noticeId, hash, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建公告失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "添加公告成功"})
}
func GetNoticeDetail(c *gin.Context) {
	type NoticeDetail struct {
		Base   dbMod.NoticeWithAuthorName `json:"base"`
		Record []fabric.ResultNotice      `json:"record"`
	}
	var detail NoticeDetail
	//获取路径id值
	id := c.Param("id")
	notice, err := dbMod.GetNoticeWithAuthorName(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取公告失败:" + err.Error()})
		return
	}
	history, err := fabric.GetNoticeHistory(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取公告历史失败:" + err.Error()})
		return
	}
	detail.Base = *notice
	detail.Record = history
	c.JSON(http.StatusOK, gin.H{"data": detail})
}
func GetNoticeAllPage(c *gin.Context) {
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
	notices, err := dbMod.GetAllNoticesWithPagination(iPage, iPageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取公告失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": notices, "total": len(notices)})
}

func GetNoticeByConditions(c *gin.Context) {
	//获取page和pageSize
	page, pageSize := c.Query("page"), c.Query("pageSize")
	if page == "" && pageSize == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "传入的page和pageSize参数不合法:"})
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
	var conditions map[string]interface{}
	err = c.ShouldBind(&conditions)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "传入的参数不合法:" + err.Error()})
		return
	}
	notices, err := dbMod.GetAllNoticesWithConditions(iPage, iPageSize, conditions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取公告失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": notices, "total": len(notices)})

}
func UpdateNotice(c *gin.Context) {
	//获取路径id值
	id := c.Param("id")
	var noticeReq dbMod.Notice
	err := c.ShouldBind(&noticeReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "传入的参数不合法:" + err.Error()})
		return
	}
	err = dbMod.UpdateNotice(id, noticeReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新公告失败:" + err.Error()})
		return
	}
	//计算hash值
	hash, err := utils.ComputeHash(noticeReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "计算公告hash失败:" + err.Error()})
		return
	}
	userId := c.MustGet("userId").(string)
	err = fabric.UpdateNotice(id, hash, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新公告失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "更新公告成功"})
}
func DeleteNotice(c *gin.Context) {
	//获取路径id值
	id := c.Param("id")
	err := dbMod.DeleteNotice(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除公告失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "删除公告成功"})
}
