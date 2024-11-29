package models

import (
	"community-governance/db"
	"fmt"
)

// VoteRule 表示投票规则表
type VoteRule struct {
	RuleID      string `gorm:"primaryKey;type:varchar(64);not null" json:"rule_id"`
	RuleType    string `gorm:"type:varchar(10);not null" json:"rule_type"`
	RuleValue   string `gorm:"type:varchar(50);not null" json:"rule_value"`
	Description string `gorm:"type:varchar(200);not null" json:"description"`
	CreateDate  string `gorm:"type:varchar(26);not null" json:"create_date"`
	RuleName    string `gorm:"type:varchar(20);not null" json:"rule_name"`
}

func (VoteRule) TableName() string {
	return "vote_rule"
}

// CreateVoteRule 创建投票规则
func CreateVoteRule(rule *VoteRule) error {
	return db.DB.Create(rule).Error
}

func GetVoteRuleById(ruleId string) (VoteRule, error) {
	var rule VoteRule
	err := db.DB.Where("rule_id = ?", ruleId).First(&rule).Error
	return rule, err
}

// UpdateVoteRule 更新投票规则
func UpdateVoteRule(ruleId string, rule interface{}) error {
	return db.DB.Model(&VoteRule{}).Where("rule_id = ?", ruleId).Updates(rule).Error
}

func DeleteVoteRule(ruleId string) error {
	result := db.DB.Where("rule_id = ?", ruleId).Delete(&VoteRule{})
	// 检查是否发生错误
	if result.Error != nil {
		return fmt.Errorf("failed to delete asset request: %v", result.Error)
	}

	// 检查影响的行数
	if result.RowsAffected == 0 {
		return fmt.Errorf("no asset request found with request_id: %s", ruleId)
	}

	// 删除成功
	return nil
}

// GetVoteRuleNames 获取所有投票规则的名字
func GetVoteRuleNames() ([]string, error) {
	var ruleTypes []string
	err := db.DB.Model(&VoteRule{}).Select("rule_name").Distinct().Scan(&ruleTypes).Error
	return ruleTypes, err
}

func GetVoteRulesAllWithPagination(page, pageSize int) ([]VoteRule, error) {
	var voteRules []VoteRule

	// 计算 OFFSET
	offset := (page - 1) * pageSize

	// 使用 Limit 和 Offset 进行分页查询
	err := db.DB.Limit(pageSize).Offset(offset).Find(&voteRules).Error
	return voteRules, err
}
func GetVoteRulesAllWithConditions(conditions map[string]interface{}, page, pageSize int) ([]VoteRule, error) {
	var voteRules []VoteRule

	// 计算 OFFSET
	offset := (page - 1) * pageSize

	// 使用 Limit 和 Offset 进行分页查询，并根据条件构建 WHERE 子句
	err := db.DB.Where(conditions).Limit(pageSize).Offset(offset).Find(&voteRules).Error
	return voteRules, err

}

// GetRuleIdByName 根据rule_name查询rule id
func GetRuleIdByName(ruleName string) (string, error) {
	var ruleId string
	err := db.DB.Model(&VoteRule{}).Where("rule_name = ?", ruleName).Select("rule_id").Find(&ruleId).Error
	return ruleId, err
}

// GetRuleTypeByRuleId 根据rule_name查询rule type
func GetRuleTypeByRuleId(ruleId string) (string, error) {
	var ruleType string
	err := db.DB.Where("rule_id = ?", ruleId).Select("rule_type").Find(&ruleType).Error
	return ruleType, err
}
