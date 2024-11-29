package models

import (
	"community-governance/db"
	"fmt"
)

// PublicFacility 公共设施信息模型
type PublicFacility struct {
	FacilityID  string `gorm:"primaryKey;type:varchar(64);not null" json:"facility_id"` // 设施ID
	Name        string `gorm:"type:varchar(20);not null" json:"name"`                   // 设施名称
	Description string `gorm:"type:varchar(200);not null" json:"description"`           // 设施说明
	Location    string `gorm:"type:varchar(100);not null" json:"location"`              // 设施位置
	Status      string `gorm:"type:varchar(64);not null" json:"status"`                 // 设施状态
	Manager     string `gorm:"type:varchar(64);not null" json:"manager"`                // 负责人
	CreateTime  string `gorm:"type:varchar(26);not null" json:"create_time"`            // 创建时间
}

func (PublicFacility) TableName() string {
	return "public_facility"
}
func CreateFacility(p *PublicFacility) error {
	return db.DB.Create(&p).Error
}

// GetFacilityByID 查询指定公共设施信息
func GetFacilityByID(id string) (*PublicFacility, error) {
	var facility PublicFacility
	err := db.DB.First(&facility, "facility_id = ?", id).Error
	return &facility, err
}

// UpdateFacility 更新公共设施信息
func UpdateFacility(id string, updatedFields interface{}) error {
	return db.DB.Model(&PublicFacility{}).Where("facility_id = ?", id).Updates(updatedFields).Error
}

// DeleteFacility 删除公共设施信息
func DeleteFacility(id string) error {
	result := db.DB.Delete(&PublicFacility{}, "facility_id = ?", id)
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

// GetAllFacilityWithPagination 获取所有公共设施并分页
func GetAllFacilityWithPagination(page, pageSize int) ([]PublicFacility, error) {
	var facility []PublicFacility
	err := db.DB.Offset((page - 1) * pageSize).Limit(pageSize).Find(&facility).Error
	return facility, err
}

// GetAllFacilityWithConditions 获取满足多条件的公共设施信息
func GetAllFacilityWithConditions(conditions map[string]interface{}, page, pageSize int) ([]PublicFacility, error) {
	var facility []PublicFacility

	// 计算 OFFSET
	offset := (page - 1) * pageSize

	// 使用 Limit 和 Offset 进行分页查询，并根据条件构建 WHERE 子句
	err := db.DB.Where(conditions).Limit(pageSize).Offset(offset).Find(&facility).Error
	return facility, err
}
