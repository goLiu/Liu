package models

import (
	"community-governance/db"
	"fmt"
)

// Notice 表对应的模型
type Notice struct {
	NoticeID    string `gorm:"primaryKey;type:varchar(64);not null" json:"notice_id"` // 公告ID
	Title       string `gorm:"type:varchar(20);not null" json:"title"`                // 公告标题
	Content     string `gorm:"type:text;not null" json:"content"`                     // 公告内容
	Type        string `gorm:"type:varchar(10);not null" json:"type"`                 // 类型
	Author      string `gorm:"type:varchar(64);not null" json:"author"`               // 发布人
	PublishTime string `gorm:"type:datetime;not null" json:"publish_time"`            // 发布时间
	Version     int    `gorm:"not null" json:"version"`                               // 发布版本
}
type NoticeWithAuthorName struct {
	NoticeID    string `json:"notice_id"`    // 公告ID
	Title       string `json:"title"`        // 公告标题
	Content     string `json:"content"`      // 公告内容
	Type        string `json:"type"`         // 类型
	AuthorName  string `json:"author_name"`  // 作者姓名
	PublishTime string `json:"publish_time"` // 发布时间
	Version     int    `json:"version"`      // 发布版本
}

func (Notice) TableName() string {
	return "notice"
}

// CreateNotice 增加公告
func CreateNotice(notice Notice) error {
	return db.DB.Create(&notice).Error
}

// GetNoticeByID 查询公告
func GetNoticeByID(id string) (*Notice, error) {
	var notice Notice
	err := db.DB.First(&notice, "notice_id = ?", id).Error
	return &notice, err
}

// UpdateNotice 更新公告
func UpdateNotice(id string, updatedFields interface{}) error {
	return db.DB.Model(&Notice{}).Where("notice_id = ?", id).Updates(updatedFields).Error
}

// DeleteNotice 删除公告
func DeleteNotice(id string) error {
	result := db.DB.Delete(&Notice{}, "notice_id = ?", id)
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

// GetAllNoticesWithPagination 获取所有公告并分页
func GetAllNoticesWithPagination(page, pageSize int) ([]Notice, error) {
	var notices []Notice
	err := db.DB.Offset((page - 1) * pageSize).Limit(pageSize).Find(&notices).Error
	return notices, err
}

func GetAllNoticesWithConditions(page, pageSize int, conditions map[string]interface{}) ([]Notice, error) {
	var notices []Notice

	// 计算 OFFSET
	offset := (page - 1) * pageSize

	// 使用 Limit 和 Offset 进行分页查询，并根据条件构建 WHERE 子句
	err := db.DB.Where(conditions).Limit(pageSize).Offset(offset).Find(&notices).Error
	return notices, err
}
func GetNoticeWithAuthorName(id string) (*NoticeWithAuthorName, error) {
	var result NoticeWithAuthorName

	// 查询公告及作者姓名
	err := db.DB.Table("notice").
		Select("notice.notice_id, notice.title, notice.content, notice.type, notice.publish_time, notice.version, member.name as author_name").
		Joins("left join member on notice.author = member.member_id").
		Where("notice.notice_id = ?", id).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}
