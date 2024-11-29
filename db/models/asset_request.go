package models

import (
	"community-governance/db"
	"fmt"
)

type AssetRequest struct {
	RequestID    string `gorm:"primaryKey;type:varchar(64);not null" json:"request_id"` // 资产申请ID
	Name         string `gorm:"type:varchar(20);not null" json:"name"`                  // 资产名称
	RequestDate  string `gorm:"type:varchar(26);not null" json:"request_date"`          // 申请时间
	Description  string `gorm:"type:varchar(200)" json:"description"`                   // 申请说明
	Status       string `gorm:"type:varchar(10);not null" json:"status"`                // 申请状态
	PurchaseTime string `gorm:"type:varchar(26);not null" json:"purchase_time"`         // 购买日期
	ProcessDate  string `gorm:"type:varchar(26);not null" json:"process_date"`          // 处理日期
	Type         string `gorm:"type:varchar(20);not null" json:"type"`                  // 资产类型
	Location     string `gorm:"type:varchar(100);not null" json:"location"`             // 资产位置
	RequestValue string `gorm:"type:varchar(100);not null" json:"request_value"`        // 申请金额
	Requester    string `gorm:"type:varchar(100);not null" json:"requester"`            // 申请人
	RequestType  string `gorm:"type:varchar(100);not null" json:"request_type"`         // 申请类型
	Asset        string `gorm:"type:varchar(64);not null" json:"asset"`                 // 资产ID
}

func (AssetRequest) TableName() string {
	return "asset_request"
}

// CreateAssetRequest 创建资产申请
func CreateAssetRequest(request AssetRequest) error {
	return db.DB.Create(&request).Error
}

// GetAssetRequestByID  查询资产申请
func GetAssetRequestByID(requestID string) (*AssetRequest, error) {
	var request AssetRequest
	err := db.DB.First(&request, "request_id = ?", requestID).Error
	return &request, err
}

// UpdateAssetRequest 更新资产申请
func UpdateAssetRequest(requestID string, updatedData interface{}) error {
	return db.DB.Model(&AssetRequest{}).Where("request_id = ?", requestID).Updates(updatedData).Error
}

// DeleteAssetRequest 删除资产申请
func DeleteAssetRequest(requestID string) error {
	// 执行删除操作
	result := db.DB.Delete(&AssetRequest{}, "request_id = ?", requestID)

	// 检查是否发生错误
	if result.Error != nil {
		return fmt.Errorf("failed to delete asset request: %v", result.Error)
	}

	// 检查影响的行数
	if result.RowsAffected == 0 {
		return fmt.Errorf("no asset request found with request_id: %s", requestID)
	}

	// 删除成功
	return nil
}

func GetAllAsseRequestWithConditions(page, pageSize int, conditions map[string]interface{}) ([]AssetRequest, error) {
	var asset []AssetRequest

	// 计算 OFFSET
	offset := (page - 1) * pageSize

	// 使用 Limit 和 Offset 进行分页查询，并根据条件构建 WHERE 子句
	err := db.DB.Where(conditions).Limit(pageSize).Offset(offset).Find(&asset).Error
	return asset, err
}

func GetAllAsseRequestWithPagination(page, pageSize int) ([]AssetRequest, error) {
	var asset []AssetRequest
	err := db.DB.Offset((page - 1) * pageSize).Limit(pageSize).Find(&asset).Error
	return asset, err
}
