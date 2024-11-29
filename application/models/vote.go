package models

type CreateVote struct {
	Name        string       `json:"name"`        //投票标题
	RuleName    string       `json:"rule_name"`   //规则名称
	StartTime   string       `json:"start_time"`  //开始时间
	Manager     string       `json:"manager"`     //负责人ID
	Description string       `json:"description"` //描述
	Options     []VoteOption `json:"options"`
}
type VoteOption struct {
	OptionValue string `json:"option_value"`
	Description string `json:"description"`
}

type CreateVoteRule struct {
	RuleType    string `json:"rule_type"`
	RuleValue   string `json:"rule_value"`
	Description string `json:"description"`
	RuleName    string `json:"rule_name"`
}
