package models

import "community-governance/db"

// VoteOption 表示投票选项表
type VoteOption struct {
	OptionID    string `gorm:"primaryKey;type:varchar(64);not null" json:"option_id"`
	VoteID      string `gorm:"type:varchar(64);not null" json:"vote_id"`
	OptionValue string `gorm:"type:varchar(100);not null" json:"option_value"`
	Status      string `gorm:"type:varchar(10);not null" json:"status"`
	CreateDate  string `gorm:"type:varchar(26);not null" json:"create_date"`
	Description string `gorm:"type:varchar(200);not null" json:"description"`
}

func (VoteOption) TableName() string {
	return "vote_option"
}

// CreateVoteOption 创建投票选项
func CreateVoteOption(option *VoteOption) error {
	return db.DB.Create(option).Error
}

// GetVoteOptionByVoteId 根据投票id获取选项
func GetVoteOptionByVoteId(voteId string) ([]VoteOption, error) {
	var options []VoteOption
	err := db.DB.Where("vote_id = ?", voteId).Find(&options).Error
	return options, err
}
