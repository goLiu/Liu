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

func AddAssetRequestRecord(c *gin.Context) {
	var assetReq models.CreateAssetRequest
	if err := c.ShouldBindJSON(&assetReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//获取userId
	userId := c.MustGet("userId").(string)
	request := dbMod.AssetRequest{
		RequestID:    uuid.New().String(),
		Name:         assetReq.Name,
		RequestDate:  utils.GetNowTimeString(),
		Description:  assetReq.Description,
		Status:       "active",
		PurchaseTime: assetReq.PurchaseTime,
		Type:         assetReq.Type,
		Location:     assetReq.Location,
		RequestValue: assetReq.RequestValue,
		RequestType:  assetReq.RequestType,
		Asset:        assetReq.Asset,
		Requester:    userId,
	}
	err := dbMod.CreateAssetRequest(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加资产申请失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "添加资产申请成功"})

}

func GetAssetRequestAllPage(c *gin.Context) {
	//获取page和pageSize
	page, pageSize := c.Query("page"), c.Query("pageSize")
	if page == "" && pageSize == "" {
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
	assetRequests, err := dbMod.GetAllAsseRequestWithPagination(iPage, iPageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取资产申请失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": assetRequests, "total": len(assetRequests)})
}

func GetAssetRequestByConditions(c *gin.Context) {
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
	assetRequests, err := dbMod.GetAllAsseRequestWithConditions(iPage, iPageSize, conditions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取资产申请失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": assetRequests, "total": len(assetRequests)})
}
func DeleteAssetRequest(c *gin.Context) {
	//获取路径id值
	id := c.Param("id")
	err := dbMod.DeleteAssetRequest(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除资产申请失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "删除资产申请成功"})
}
func AuditRequestRecord(c *gin.Context) {
	//获取路径id值
	id := c.Param("id")
	status := c.Query("status")
	request, err := dbMod.GetAssetRequestByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取资产申请失败:" + err.Error()})
		return
	}
	request.Status = status
	request.PurchaseTime = utils.GetNowTimeString()

	userId := c.MustGet("userId").(string)
	//判断状态
	if status == "pass" {
		//判断申请类型
		if request.RequestType == "add" {
			asset := dbMod.Asset{
				AssetID:      id,
				Name:         request.Name,
				Type:         request.Type,
				Description:  request.Description,
				Status:       "pending",
				Location:     request.Location,
				PurchaseDate: request.PurchaseTime,
				Owner:        request.Requester,
			}
			//计算hash值
			hash, err := utils.ComputeHash(asset)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "计算资产hash失败:" + err.Error()})
				return
			}
			//将hash值写入区块链
			err = fabric.CreateAsset(asset.AssetID, hash, asset.Owner, userId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "更新资产失败:" + err.Error()})
				return
			}
			err = dbMod.CreateAsset(&asset)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "创建资产失败:" + err.Error()})
				return
			}
		} else if request.RequestType == "update" {
			asset, err := dbMod.GetAssetByID(request.Asset)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "获取资产失败:" + err.Error()})
				return
			}
			asset.Owner = request.RequestValue
			//将hash值写入区块链
			err = fabric.ExchangeOwner(asset.AssetID, asset.Owner)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "更新资产失败:" + err.Error()})
				return
			}
			err = dbMod.UpdateAsset(asset.AssetID, asset)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "更新资产失败:" + err.Error()})
				return
			}

		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "申请类型不合法"})
			return
		}
	} else if status != "reject" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "状态不合法"})
		return
	}
	err = dbMod.UpdateAssetRequest(id, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "审核资产申请失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "审核资产申请成功"})
}
func GetAssetRequestByPerson(c *gin.Context) {
	//获取userId
	userId := c.MustGet("userId").(string)
	//获取page 和 pageSize
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

	condition := make(map[string]interface{})
	condition["requester"] = userId
	assetRequests, err := dbMod.GetAllAsseRequestWithConditions(iPage, iPageSize, condition)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取资产申请失败:" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": assetRequests, "total": len(assetRequests)})
}

//func GetAssetRequestDetail(c *gin.Context) {
//
//}
