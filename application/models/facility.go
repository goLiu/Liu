package models

type CreateFacility struct {
	Name        string `json:"name"`        // 设施名称
	Description string `json:"description"` // 设施说明
	Location    string `json:"location"`    // 设施位置
	Manager     string `json:"manager"`     // 负责人
}
