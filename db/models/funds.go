package models

import (
	"community-governance/db"
	"fmt"
)

// Fund 表对应的模型
type Fund struct {
	FundID            string `gorm:"primaryKey;type:varchar(64);not null" json:"fund_id"` // 款项ID
	Description       string `gorm:"type:varchar(200);not null" json:"description"`       // 款项说明
	Status            string `gorm:"type:varchar(10);not null" json:"status"`             // 款项状态
	Name              string `gorm:"type:varchar(20);not null" json:"name"`               // 款项名称
	Source            string `gorm:"type:varchar(10);not null" json:"source"`             // 款项来源
	SourceDescription string `gorm:"type:varchar(100)" json:"source_description"`         // 款项来源说明
	TotalAmount       string `gorm:"type:varchar(20);not null" json:"total_amount"`       // 款项的总金额
	CurrentBalance    string `gorm:"type:varchar(20);not null" json:"current_balance"`    // 当前款项余额
	Manager           string `gorm:"type:varchar(64);not null" json:"manager"`            // 负责人
	EstablishDate     string `gorm:"type:varchar(26);not null" json:"establish_date"`     // 成立时间
}

func (f *Fund) TableName() string {
	return "fund"
}
func CreateFunds(funds *Fund) error {
	return db.DB.Create(&funds).Error
}

// GetFundByID 查询公告
func GetFundByID(id string) (*Fund, error) {
	var Fund Fund
	err := db.DB.First(&Fund, "fund_id = ?", id).Error
	return &Fund, err
}

// UpdateFund 更新公告
func UpdateFund(id string, updatedFields interface{}) error {
	return db.DB.Model(&Fund{}).Where("fund_id = ?", id).Updates(updatedFields).Error
}

// DeleteFund 删除公告
func DeleteFund(id string) error {

	result := db.DB.Delete(&Fund{}, "fund_id = ?", id)
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

// GetAllFundsWithPagination 获取所有公告并分页
func GetAllFundsWithPagination(page, pageSize int) ([]Fund, error) {
	var fund []Fund
	err := db.DB.Offset((page - 1) * pageSize).Limit(pageSize).Find(&fund).Error
	return fund, err
}
func GetAllFundsWithConditions(conditions map[string]interface{}, page, pageSize int) ([]Fund, error) {
	var fund []Fund

	// 计算 OFFSET
	offset := (page - 1) * pageSize

	// 使用 Limit 和 Offset 进行分页查询，并根据条件构建 WHERE 子句
	err := db.DB.Where(conditions).Limit(pageSize).Offset(offset).Find(&fund).Error
	return fund, err
}
