package models

import (
	"community-governance/db"
	"fmt"
)

// Vote 表示投票表
type Vote struct {
	VoteID      string `gorm:"primaryKey;type:varchar(64);not null" json:"vote_id"`
	Name        string `gorm:"type:varchar(20);" json:"name"`
	RuleID      string `gorm:"type:varchar(64);not null" json:"rule_id"`
	StartTime   string `gorm:"type:varchar(26);not null" json:"start_time"`
	Manager     string `gorm:"type:varchar(64);not null" json:"manager"`
	Description string `gorm:"type:varchar(200);" json:"description"`
	Status      string `gorm:"type:varchar(10);not null" json:"status"`
	Result      string `gorm:"type:varchar(64);" json:"result"`
}

func (Vote) TableName() string {
	return "vote"
}

// CreateVote 创建投票
func CreateVote(vote *Vote) error {
	return db.DB.Create(vote).Error
}

// GetVoteByID 获取指定投票信息
func GetVoteByID(memberId string) (*Vote, error) {
	var vote Vote
	err := db.DB.Where("vote_id = ?", memberId).First(&vote).Error
	return &vote, err
}

// GetVoteAllWithPagination 获取分页后的投票信息
func GetVoteAllWithPagination(page, pageSize int) ([]Vote, error) {
	var votes []Vote

	// 计算 OFFSET
	offset := (page - 1) * pageSize

	// 使用 Limit 和 Offset 进行分页查询
	err := db.DB.Limit(pageSize).Offset(offset).Find(&votes).Error
	return votes, err
}

// UpdateVote 更新指定成员信息
func UpdateVote(voteId, vote interface{}) error {
	return db.DB.Model(&Vote{}).Where("vote_id = ?", voteId).Updates(vote).Error
}

// DeleteVote 删除指定成员信息
func DeleteVote(voteId string) error {
	result := db.DB.Where("vote_id = ?", voteId).Delete(&Vote{})
	// 检查是否发生错误
	if result.Error != nil {
		return fmt.Errorf("failed to delete asset request: %v", result.Error)
	}

	// 检查影响的行数
	if result.RowsAffected == 0 {
		return fmt.Errorf("no asset request found with request_id: %s", voteId)
	}

	// 删除成功
	return nil
}

func GetVoteAllWithConditions(conditions map[string]interface{}, page, pageSize int) ([]Vote, error) {
	var votes []Vote

	// 计算 OFFSET
	offset := (page - 1) * pageSize

	// 使用 Limit 和 Offset 进行分页查询，并根据条件构建 WHERE 子句
	err := db.DB.Where(conditions).Limit(pageSize).Offset(offset).Find(&votes).Error
	return votes, err
}
