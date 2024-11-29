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

const (
	statusInit = "use"
)

func GetFacilityDetail(c *gin.Context) {
	type Info struct {
		Base    dbMod.PublicFacility `json:"base"`
		Records []fabric.UsageRecord `json:"records"`
	}
	//获取路径id值
	id := c.Param("id")
	facility, err := dbMod.GetFacilityByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取设施失败:" + err.Error()})
		return
	}
	history, err := fabric.GetFacilityUsageHistory(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取设施记录失败:" + err.Error()})
		return
	}

	info := Info{
		Base:    *facility,
		Records: history,
	}
	c.JSON(http.StatusOK, gin.H{"data": info})

}
func AddFacility(c *gin.Context) {
	var facilityReq models.CreateFacility
	err := c.ShouldBind(&facilityReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "传入的参数不合法:" + err.Error()})
		return
	}
	facilityId := uuid.New().String()
	facility := dbMod.PublicFacility{
		FacilityID:  facilityId,
		Name:        facilityReq.Name,
		Description: facilityReq.Description,
		Location:    facilityReq.Location,
		Status:      statusInit,
		Manager:     facilityReq.Manager,
		CreateTime:  utils.GetNowTimeString(),
	}
	err = dbMod.CreateFacility(&facility)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加设施失败:" + err.Error()})
		return
	}
	//计算hash
	hash, err := utils.ComputeHash(facility)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "计算设施hash失败:" + err.Error()})
		return
	}
	err = fabric.RegisterFacility(facilityId, hash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加设施失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "添加设施成功"})
}

func GetFacilityAllPage(c *gin.Context) {
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
	facilities, err := dbMod.GetAllFacilityWithPagination(iPage, iPageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取设施失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": facilities, "total": len(facilities)})
}
func UpdateFacility(c *gin.Context) {
	var facilityReq dbMod.PublicFacility
	err := c.ShouldBind(&facilityReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "传入的参数不合法:" + err.Error()})
		return
	}
	facilityReq.CreateTime = utils.GetNowTimeString()
	err = dbMod.UpdateFacility(facilityReq.FacilityID, facilityReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新设施失败:" + err.Error()})
		return
	}
	//计算hash
	hash, err := utils.ComputeHash(facilityReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "计算设施hash失败:" + err.Error()})
		return
	}
	err = fabric.UpdateFacility(facilityReq.FacilityID, hash, facilityReq.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新设施失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "更新设施成功"})
}
func DeleteFacility(c *gin.Context) {
	//获取路径id值
	id := c.Param("id")
	err := dbMod.DeleteFacility(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除设施失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "删除设施成功"})

}
func GetFacilityByConditions(c *gin.Context) {
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
	facilities, err := dbMod.GetAllFacilityWithConditions(conditions, iPage, iPageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取设施失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": facilities, "total": len(facilities)})
}
func RequestFacility(c *gin.Context) {
	//获取路径id值
	id := c.Param("id")
	//获取userId
	userId := c.MustGet("userId").(string)
	err := fabric.RequestFacility(id, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "请求设施失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "请求设施成功"})

}
func ReleaseFacility(c *gin.Context) {
	//获取路径id值
	id := c.Param("id")
	//获取userId
	userId := c.MustGet("userId").(string)
	err := fabric.ReleaseFacility(id, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "释放设施失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "释放设施成功"})
}
