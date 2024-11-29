package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	"sort"
	"strconv"
	"strings"
)

const (
	RuleTypeMajority  = "majority"  //简单多数法
	RuleTypeThreshold = "threshold" //阈值法
	RuleTypeTopN      = "top_n"
)

type VoteContract struct {
	contractapi.Contract
}

type Vote struct {
	BaseHash  string         `json:"base_hash"`  //投票基础信息
	RuleType  string         `json:"rule_type"`  //投票规则
	RuleValue string         `json:"rule_value"` //投票规则涉及值
	Options   map[string]int `json:"options"`    //投票选项与票数
	IsEnd     bool           `json:"is_end"`     //是否结束
}
type VoteRecord struct {
	Voter    string `json:"voter"`     //投票人
	Option   string `json:"option"`    //投票选项
	VoteTime string `json:"vote_time"` //投票时间
}

func (v *VoteContract) CreatVote(ctx contractapi.TransactionContextInterface, id, base, ruleType, ruleValue, options string) error {
	state, err := ctx.GetStub().GetState(id)
	if err != nil {
		return err
	}
	if state != nil {
		return fmt.Errorf("%s is exist", id)
	}
	//将options转换为map
	optionSl := strings.Split(options, ",")
	optionMap := make(map[string]int)
	for _, opt := range optionSl {
		optionMap[opt] = 0
	}
	vote := Vote{
		BaseHash:  base,
		RuleType:  ruleType,
		RuleValue: ruleValue,
		Options:   optionMap,
		IsEnd:     false,
	}
	data, err := json.Marshal(vote)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(id, data)
}

// VoteJoin 投票
func (v *VoteContract) VoteJoin(ctx contractapi.TransactionContextInterface, id, voter, option string) (string, error) {
	//获取投票信息
	vote, err := v.GetVote(ctx, id)
	if err != nil {
		return "", err
	}
	//判断投票项目是否已经结束
	if vote.IsEnd {
		return "", fmt.Errorf("vote is end")
	}
	result := ""
	//更新选项票数
	vote.Options[option]++
	if vote.RuleType == RuleTypeThreshold {
		iVal, err := strconv.Atoi(vote.RuleValue)
		if err != nil {
			return "", fmt.Errorf("invalid rule value:%s", err.Error())
		}
		if vote.Options[option] >= iVal {
			vote.IsEnd = true
		}
		result = option
	}
	data, err := json.Marshal(vote)
	if err != nil {
		return "", fmt.Errorf("failed marshal:%s", err.Error())
	}
	err = ctx.GetStub().PutState(id, data)
	if err != nil {
		return "", fmt.Errorf("failed to put vote state:%s", err.Error())
	}
	notTime, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return "", fmt.Errorf("failed to get tx timestamp:%s", err.Error())
	}
	//更新投票记录
	voteRecord := VoteRecord{
		Voter:    voter,
		Option:   option,
		VoteTime: notTime.AsTime().Format("2006-01-02 15:04:05"),
	}
	data, err = json.Marshal(voteRecord)
	if err != nil {
		return "", err
	}
	recordId := fmt.Sprintf("%s-%s", "record", id)
	err = ctx.GetStub().PutState(recordId, data)
	if err != nil {
		return "", fmt.Errorf("failed to put vote record state:%s", err.Error())
	}

	return result, nil
}

func (v *VoteContract) GetVote(ctx contractapi.TransactionContextInterface, id string) (Vote, error) {
	state, err := ctx.GetStub().GetState(id)
	if err != nil {
		return Vote{}, err
	}
	if state == nil {
		return Vote{}, fmt.Errorf("%s is not exist", id)
	}

	var vote Vote
	if err := json.Unmarshal(state, &vote); err != nil {
		return Vote{}, err
	}
	return vote, nil
}

func (v *VoteContract) GetVoteRecordHistory(ctx contractapi.TransactionContextInterface, id string) ([]VoteRecord, error) {
	recordId := fmt.Sprintf("%s-%s", "record", id)
	//判断是否存在
	state, err := ctx.GetStub().GetState(recordId)
	if err != nil {
		return nil, err
	}
	if state == nil {
		return nil, nil
	}
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(recordId)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	var records []VoteRecord
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var record VoteRecord
		if err := json.Unmarshal(queryResult.Value, &record); err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}

// EndVote 结束投票，获取投票结果
func (v *VoteContract) EndVote(ctx contractapi.TransactionContextInterface, id string) ([]string, error) {
	vote, err := v.GetVote(ctx, id)
	if err != nil {
		return nil, err
	}
	if vote.IsEnd {
		return nil, fmt.Errorf("vote is end")
	}
	result := make([]string, 0)
	//获取type
	ruleType := vote.RuleType
	if ruleType == RuleTypeMajority {
		r := processMajority(vote.Options)
		result = append(result, r)
	} else if ruleType == RuleTypeTopN {
		iVal, err := strconv.Atoi(vote.RuleValue)
		if err != nil {
			return nil, fmt.Errorf("invalid rule value:%s", err.Error())
		}
		result = processTopN(vote.Options, iVal)
	} else {
		return nil, fmt.Errorf("invalid rule type:%s,only support majority and topN", ruleType)
	}
	return result, nil
}
func processMajority(options map[string]int) string {
	//获取options 中数字最高的key
	maxKey := ""
	maxValue := 0
	for key, value := range options {
		if value > maxValue {
			maxKey = key
			maxValue = value
		}
	}
	return maxKey
}
func processTopN(options map[string]int, n int) []string {
	// 存储键值对
	type kv struct {
		Key   string
		Value int
	}

	var sortedOptions []kv
	for key, value := range options {
		sortedOptions = append(sortedOptions, kv{Key: key, Value: value})
	}

	// 按值从大到小排序
	sort.Slice(sortedOptions, func(i, j int) bool {
		return sortedOptions[i].Value > sortedOptions[j].Value
	})

	// 取前 n 个键
	var result []string
	for i := 0; i < n && i < len(sortedOptions); i++ {
		result = append(result, sortedOptions[i].Key)
	}

	return result
}

// CloseVote 异常管理投票
func (v *VoteContract) CloseVote(ctx contractapi.TransactionContextInterface, id string) error {
	vote, err := v.GetVote(ctx, id)
	if err != nil {
		return err
	}
	vote.IsEnd = true
	data, err := json.Marshal(vote)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(id, data)
}
