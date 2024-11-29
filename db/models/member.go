package models

import (
	"community-governance/db"
	"errors"
	"fmt"
)

// Member 社区成员表
type Member struct {
	MemberID      string `gorm:"column:member_id;primaryKey;type:varchar(64);not null" json:"member_id"`
	Name          string `gorm:"column:name;type:varchar(16);not null" json:"name"`
	Type          string `gorm:"column:type;type:varchar(10);not null" json:"type"`
	HouseholdID   string `gorm:"column:household_id;type:varchar(64);not null" json:"household_id"`
	IDNumber      string `gorm:"column:id_number;type:varchar(18);not null" json:"id_number"`
	Address       string `gorm:"column:address;type:varchar(100);not null" json:"address"`
	Sex           string `gorm:"column:sex;type:varchar(10);not null" json:"sex"`
	DateBirth     string `gorm:"column:date_birth;type:varchar(26);not null" json:"date_birth"`
	State         string `gorm:"column:state;type:varchar(10);not null" json:"state"`
	Phone         string `gorm:"column:phone;type:varchar(20);not null" json:"phone"`
	NameUsed      string `gorm:"column:name_used;type:varchar(16)" json:"name_used"`
	Remarks       string `gorm:"column:remarks;type:varchar(100)" json:"remarks"`
	Education     string `gorm:"column:education;type:varchar(10);not null" json:"education"`
	MaritalStatus string `gorm:"column:marital_status;type:varchar(10);not null" json:"marital_status"`
	Password      string `gorm:"column:password;type:varchar(20);not null" json:"-"`
}

func (Member) TableName() string {
	return "member"
}

// CreateMember 新增一名新成员
func CreateMember(member *Member) error {
	return db.DB.Create(member).Error
}

// GetMemberByID 获取指定成员信息
func GetMemberByID(memberID string) (*Member, error) {
	var member Member
	err := db.DB.Where("member_id = ?", memberID).First(&member).Error
	return &member, err
}

// GetAllMembersWithPagination 获取分页后的成员信息
func GetAllMembersWithPagination(page, pageSize int) ([]Member, error) {
	var members []Member

	// 计算 OFFSET
	offset := (page - 1) * pageSize

	// 使用 Limit 和 Offset 进行分页查询
	err := db.DB.Limit(pageSize).Offset(offset).Find(&members).Error
	return members, err
}

// UpdateMemberAttribute 更新指定成员信息的指定字段信息
//func UpdateMemberAttribute(memberID string, updatedData map[string]interface{}) error {
//	return db.DB.Model(&Member{}).Where("member_id = ?", memberID).Updates(updatedData).Error
//}

// UpdateMember 更新指定成员信息
func UpdateMember(memberID string, member interface{}) error {
	return db.DB.Model(&Member{}).Where("member_id = ?", memberID).Updates(member).Error
}

// DeleteMember 删除指定成员信息
func DeleteMember(memberID string) error {
	result := db.DB.Where("member_id = ?", memberID).Delete(&Member{})
	// 检查是否发生错误
	if result.Error != nil {
		return fmt.Errorf("failed to delete asset request: %v", result.Error)
	}

	// 检查影响的行数
	if result.RowsAffected == 0 {
		return fmt.Errorf("no asset request found with request_id: %s", memberID)
	}

	// 删除成功
	return nil
}

// GetAllMembersWithConditions 获取满足多条件的成员信息
func GetAllMembersWithConditions(conditions map[string]interface{}, page, pageSize int) ([]Member, error) {
	var members []Member

	// 计算 OFFSET
	offset := (page - 1) * pageSize

	// 使用 Limit 和 Offset 进行分页查询，并根据条件构建 WHERE 子句
	err := db.DB.Where(conditions).Limit(pageSize).Offset(offset).Find(&members).Error
	return members, err
}

// Login 根据身份证号码和密码进行用户验证，返回用户的 ID
func Login(idNumber, password string) (string, string, error) {
	var member Member

	// 根据身份证号码查找用户
	err := db.DB.Where("id_number = ? AND password = ?", idNumber, password).First(&member).Error
	if err != nil {
		return "", "", errors.New("invalid id number or password")
	}

	// 返回ID
	return member.MemberID, member.Type, nil
}

// GetMemberNameByID 根据ID获取成员Name
func GetMemberNameByID(memberID string) (string, error) {
	var member Member
	err := db.DB.Model(&Member{}).Where("member_id = ?", memberID).First(&member).Error
	if err != nil {
		return "", err
	}
	return member.Name, nil
}
