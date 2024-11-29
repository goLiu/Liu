package models

type CreateAsset struct {
	Name         string `json:"name"`          // 资产名称
	Type         string `json:"type"`          // 资产类型
	Description  string `json:"description"`   // 资产说明
	Location     string `json:"location"`      // 资产位置
	PurchaseDate string `json:"purchase_date"` // 购买日期
	Owner        string `json:"owner"`         // 资产拥有者
}

type CreateAssetRequest struct {
	Name         string `json:"name"`           // 资产名称
	Description  string ` json:"description"`   // 申请说明
	Status       string `json:"status"`         // 申请状态
	PurchaseTime string `json:"purchase_time"`  // 购买日期
	ProcessDate  string `json:"process_date"`   // 处理日期
	Type         string `json:"type"`           // 资产类型
	Location     string `json:"location"`       // 资产位置
	RequestValue string ` json:"request_value"` // 申请传递得值
	RequestType  string `json:"request_type"`   // 申请类型
	Asset        string `json:"asset"`          // 资产
}
