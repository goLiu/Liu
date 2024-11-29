package models

type CreateFund struct {
	Description       string ` json:"description"`        // 款项说明
	Name              string `json:"name"`                // 款项名称
	Source            string ` json:"source"`             // 款项来源
	SourceDescription string ` json:"source_description"` // 款项来源说明
	TotalAmount       string ` json:"total_amount"`       // 款项的总金额
	Manager           string ` json:"manager"`            // 负责人
	EstablishDate     string `json:"establish_date"`      // 成立时间

}
type FinancialRecord struct {
	Type    string `json:"type"`    //记录类型，0-收入 1-支出
	Amount  string `json:"amount"`  //记录金额
	Source  string `json:"source"`  //金额来源
	Explain string `json:"explain"` //说明
}
