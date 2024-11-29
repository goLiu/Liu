package models

type CreateMember struct {
	Name          string `json:"name" binding:"required"`             // 姓名
	HouseholdID   string `json:"household_id" binding:"required"`     // 户号
	IDNumber      string `json:"id_number" binding:"required,len=18"` // 身份证号码，要求18位
	Address       string `json:"address" binding:"required"`          // 居住地址
	Sex           string `json:"sex" binding:"required"`              // 性别
	DateBirth     string `json:"date_birth" binding:"required"`       // 出生日期
	Phone         string `json:"phone" binding:"required"`            // 电话号码
	NameUsed      string `json:"name_used"`                           // 曾用名
	Remarks       string `json:"remarks"`                             // 备注
	Education     string `json:"education" binding:"required"`        // 学历
	MaritalStatus string `json:"marital_status" binding:"required"`   // 婚姻状况
	Type          string `json:"type" binding:"required"`             // 类型
}
