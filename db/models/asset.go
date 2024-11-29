package models

import (
	"community-governance/db"
	"fmt"
)

type Asset struct {
	AssetID      string `gorm:"primaryKey;type:varchar(64);not null" json:"asset_id"` // 资产ID
	Name         string `gorm:"type:varchar(20);not null" json:"name"`                // 资产名称
	Type         string `gorm:"type:varchar(20);not null" json:"type"`                // 资产类型
	Description  string `gorm:"type:varchar(200)" json:"description"`                 // 资产说明
	Status       string `gorm:"type:varchar(10);not null" json:"status"`              // 资产状态
	Location     string `gorm:"type:varchar(100);not null" json:"location"`           // 资产位置
	PurchaseDate string `gorm:"type:varchar(26);not null" json:"purchase_date"`       // 购买日期
	Owner        string `gorm:"type:varchar(64);not null" json:"owner"`               // 资产拥有者
}

func (Asset) TableName() string {
	return "asset"
}

func GetAssetByID(id string) (*Asset, error) {
	var asset Asset
	err := db.DB.First(&asset, "asset_id = ?", id).Error
	return &asset, err
}

func UpdateAsset(id string, updatedFields interface{}) error {
	return db.DB.Model(&Asset{}).Where("asset_id = ?", id).Updates(updatedFields).Error
}

func DeleteAsset(id string) error {
	result := db.DB.Delete(&Asset{}, "asset_id = ?", id)
	// 检查是否发生错误
	if result.Error != nil {
		return fmt.Errorf("failed to delete asset request: %v", result.Error)
	}

	// 检查影响的行数
	if result.RowsAffected == 0 {
		return fmt.Errorf("no asset request found with request_id: %s", id)
	}

	// 删除成功
	return nil
}

func GetAllAssetWithPagination(page, pageSize int) ([]Asset, error) {
	var asset []Asset
	err := db.DB.Offset((page - 1) * pageSize).Limit(pageSize).Find(&asset).Error
	return asset, err
}

func GetAllAssetWithConditions(conditions map[string]interface{}, page, pageSize int) ([]Asset, error) {
	var asset []Asset

	// 计算 OFFSET
	offset := (page - 1) * pageSize

	// 使用 Limit 和 Offset 进行分页查询，并根据条件构建 WHERE 子句
	err := db.DB.Where(conditions).Limit(pageSize).Offset(offset).Find(&asset).Error
	return asset, err
}

func CreateAsset(asset *Asset) error {
	return db.DB.Create(asset).Error
}
