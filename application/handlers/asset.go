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
	assetStatusInit = "active"
)

func AddAsset(c *gin.Context) {
	var assetReq models.CreateAsset
	if err := c.ShouldBindJSON(&assetReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	assetId := uuid.New().String()
	asset := dbMod.Asset{
		AssetID:      assetId,
		Name:         assetReq.Name,
		Type:         assetReq.Type,
		Description:  assetReq.Description,
		Status:       assetStatusInit,
		Location:     assetReq.Location,
		PurchaseDate: assetReq.PurchaseDate,
		Owner:        assetReq.Owner,
	}
	err := dbMod.CreateAsset(&asset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//计算hash值
	hash, err := utils.ComputeHash(asset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "计算资产hash失败:" + err.Error()})
		return
	}
	//获取userId
	userId := c.MustGet("userId").(string)
	err = fabric.CreateAsset(assetId, hash, assetReq.Owner, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"asset_id": assetId})
}

func GetAssetAllPage(c *gin.Context) {
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
	assets, err := dbMod.GetAllAssetWithPagination(iPage, iPageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取资产失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": assets, "total": len(assets)})
}
func UpdateAsset(c *gin.Context) {
	//获取路径id值
	var assetReq dbMod.Asset
	if err := c.ShouldBindJSON(&assetReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := dbMod.UpdateAsset(assetReq.AssetID, &assetReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//计算hash
	hash, err := utils.ComputeHash(assetReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "计算资产hash失败:" + err.Error()})
		return
	}
	//获取userId
	userId := c.MustGet("userId").(string)
	err = fabric.UpdateAsset(assetReq.AssetID, hash, assetReq.Owner, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "更新资产成功"})

}
func DeleteAsset(c *gin.Context) {
	//获取路径id值
	id := c.Param("id")
	err := dbMod.DeleteAsset(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除资产失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "删除资产成功"})

}
func GetAssetByConditions(c *gin.Context) {
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
	assets, err := dbMod.GetAllAssetWithConditions(conditions, iPage, iPageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取资产失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": assets, "total": len(assets)})
}
func GetAssetDetail(c *gin.Context) {
	type Info struct {
		Base    dbMod.Asset    `json:"base"`
		Records []fabric.Asset `json:"records"`
	}
	//获取路径id值
	id := c.Param("id")
	asset, err := dbMod.GetAssetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取资产失败:" + err.Error()})
		return
	}

	records, err := fabric.GetAssetHistory(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取资产记录失败:" + err.Error()})
		return
	}
	info := Info{
		Base:    *asset,
		Records: records,
	}
	c.JSON(http.StatusOK, gin.H{"data": info, "total": len(records)})

}
