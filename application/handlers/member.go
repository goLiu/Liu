package handlers

import (
	"community-governance/application/models"
	dbMod "community-governance/db/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

const (
	MemberStateActive   = "active"
	MemberStateInactive = "inactive"
)

// AddMember 添加成员
func AddMember(c *gin.Context) {
	var memberReq models.CreateMember
	//绑定参数
	err := c.ShouldBind(&memberReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "传入的参数不合法:" + err.Error()})
		return
	}
	//获取初始密码
	idNum := memberReq.IDNumber
	pwd := idNum[len(idNum)-6:]
	member := dbMod.Member{
		MemberID:      uuid.New().String(),
		Name:          memberReq.Name,
		HouseholdID:   memberReq.HouseholdID,
		IDNumber:      memberReq.IDNumber,
		Address:       memberReq.Address,
		Sex:           memberReq.Sex,
		DateBirth:     memberReq.DateBirth,
		State:         MemberStateActive,
		Phone:         memberReq.Phone,
		NameUsed:      memberReq.NameUsed,
		Remarks:       memberReq.Remarks,
		Education:     memberReq.Education,
		Type:          memberReq.Type,
		Password:      pwd,
		MaritalStatus: memberReq.MaritalStatus,
	}
	err = dbMod.CreateMember(&member)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加成员失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "添加成员成功"})
}

func GetMember(c *gin.Context) {
	//获取路径id值
	id := c.Param("id")
	member, err := dbMod.GetMemberByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取成员失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": member})
}

// GetMemberAllPage 获取分页获取所有成员信息
func GetMemberAllPage(c *gin.Context) {
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
	members, err := dbMod.GetAllMembersWithPagination(iPage, iPageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取成员失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": members, "total": len(members)})
}

// UpdateMember 更新成员
func UpdateMember(c *gin.Context) {
	var member dbMod.Member
	err := c.ShouldBind(&member)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "传入的参数不合法:" + err.Error()})
		return
	}
	err = dbMod.UpdateMember(member.MemberID, &member)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新成员失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "更新成员成功"})
}

// DeleteMember 删除成员
func DeleteMember(c *gin.Context) {
	//获取路径id值
	id := c.Param("id")
	err := dbMod.DeleteMember(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除成员失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "删除成员成功"})
}

// UpdateMemberState 更新成员状态
func UpdateMemberState(c *gin.Context) {
	//获取路径id值
	id := c.Param("id")
	state := c.Query("state")
	updateMap := map[string]interface{}{
		"state": state,
	}
	err := dbMod.UpdateMember(id, updateMap)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新成员状态失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "更新成员状态成功"})
}

// GetMemberByConditions 获取基于条件获取成员信息
func GetMemberByConditions(c *gin.Context) {
	//获取page和pageSize
	page, pageSize := c.Query("page"), c.Query("pageSize")
	if page == "" || pageSize == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "传入的page和pageSize参数不合法"})
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
	members, err := dbMod.GetAllMembersWithConditions(conditions, iPage, iPageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取成员失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": members, "total": len(members)})
}
