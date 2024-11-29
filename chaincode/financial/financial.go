package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// FinancialContract 财务管理合约
type FinancialContract struct {
	contractapi.Contract
}
type Financial struct {
	MessageHash string `json:"message"`     //信息hash值
	UpdateDate  string `json:"update_date"` //上次修改时间
	State       string `json:"state"`       //款项状态
}
type FinancialRecord struct {
	Type       string `json:"type"`      //记录类型，0-收入 1-支出
	Amount     string `json:"amount"`    //记录金额
	Source     string `json:"source"`    //金额来源
	Explain    string `json:"explain"`   //说明
	RecorderID string `json:"record_id"` //记录人ID
}

const (
	typeIn     = "0" //记录类型-收入
	typeOut    = "1" //记录类型-支出
	stateInit  = "0" //项目状态-运行中
	stateClose = "1" //项目已关闭
)

// CreateFinancial  创建资产项目
func (f *FinancialContract) CreateFinancial(ctx contractapi.TransactionContextInterface, id, hash string) error {
	if id == "" || hash == "" {
		return fmt.Errorf("id = ''  or hash = ''")
	}
	//判断id是否已经存在
	exist, err := f.IsExist(ctx, id)
	if err != nil {
		return err
	}
	if exist {
		return fmt.Errorf("%s already existed", id)
	}
	notTime, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("failed to get tx timestamp:%s", err.Error())
	}
	fin := Financial{
		MessageHash: hash,
		UpdateDate:  notTime.AsTime().Format("2006-01-02 15:04:05"),
		State:       stateInit,
	}
	marshal, err := json.Marshal(fin)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(id, marshal)
}
func (f *FinancialContract) UpdateFinancial(ctx contractapi.TransactionContextInterface, id, hash string) error {
	if id == "" || hash == "" {
		return fmt.Errorf("id = '' or  hash = ''")
	}
	financial, err := f.GetFinancial(ctx, id)
	if err != nil {
		return err
	}
	//判断项目状态，如果是关闭状态，则不能再执行修改的操作
	if financial.State == stateClose {
		return fmt.Errorf("%s is close", id)
	}
	financial.MessageHash = hash
	notTime, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("failed to get tx timestamp:%s", err.Error())
	}
	financial.UpdateDate = notTime.AsTime().Format("2006-01-02 15:04:05")
	marshal, err := json.Marshal(financial)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(id, marshal)

}

// AddRecord 添加款项变动记录
func (f *FinancialContract) AddRecord(ctx contractapi.TransactionContextInterface, id, finType, source, explain, recorder, amount string) error {
	recordId := fmt.Sprintf("%s-%s", "record", id)
	financial, err := f.GetFinancial(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get financial:%s", err.Error())
	}
	//判断项目状态，如果是关闭状态，则不能再执行修改的操作
	if financial.State == stateClose {
		return fmt.Errorf("%s is close", id)
	}
	record := FinancialRecord{
		Type:       finType,
		Source:     source,
		Amount:     amount,
		Explain:    explain,
		RecorderID: recorder,
	}
	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal:%s", err.Error())
	}
	return ctx.GetStub().PutState(recordId, data)
}

// ExchangeState 更改资产项目状态
func (f *FinancialContract) ExchangeState(ctx contractapi.TransactionContextInterface, id, state string) error {
	financial, err := f.GetFinancial(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get financial:%s", err.Error())
	}
	if state != stateClose && state != stateInit {
		return fmt.Errorf("illegal state value")
	}
	financial.State = state
	data, err := json.Marshal(financial)
	if err != nil {
		return fmt.Errorf("failed marshal:%s", err.Error())
	}
	return ctx.GetStub().PutState(id, data)
}

// IsExist 判断指定id是否已经存在
func (f *FinancialContract) IsExist(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	if id == "" {
		return false, fmt.Errorf("id = ''")
	}
	state, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, err
	}
	return state != nil, nil
}

// GetFinancial 查询指定款项信息
func (f *FinancialContract) GetFinancial(ctx contractapi.TransactionContextInterface, id string) (Financial, error) {
	if id == "" {
		return Financial{}, fmt.Errorf("id = ''")
	}
	state, err := ctx.GetStub().GetState(id)
	if err != nil {
		return Financial{}, err
	}
	if state == nil {
		return Financial{}, fmt.Errorf("%s not exist", id)
	}
	var fin Financial
	err = json.Unmarshal(state, &fin)
	if err != nil {
		return Financial{}, err
	}
	return fin, nil
}

func (f *FinancialContract) GetFinancialHistory(ctx contractapi.TransactionContextInterface, id string) ([]FinancialRecord, error) {
	//判断是否存在
	state, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get state:%s", err.Error())
	}
	if state == nil {
		return nil, fmt.Errorf("%s not exist", id)
	}
	records := make([]FinancialRecord, 0)
	iter, err := ctx.GetStub().GetHistoryForKey(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get history:%s", err.Error())
	}
	defer iter.Close()
	for iter.HasNext() {
		record, err := iter.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get history:%s", err.Error())
		}
		var r FinancialRecord
		if err := json.Unmarshal(record.Value, &r); err != nil {
			return nil, fmt.Errorf("failed to unmarshal:%s", err.Error())
		}
		records = append(records, r)
	}
	return records, nil
}

func (f *FinancialContract) GetFinancialRecordHistory(ctx contractapi.TransactionContextInterface, id string) ([]FinancialRecord, error) {
	recordId := fmt.Sprintf("%s-%s", "record", id)
	//判断是否存在
	state, err := ctx.GetStub().GetState(recordId)
	if err != nil {
		return nil, err
	}
	if state == nil {
		return nil, fmt.Errorf("%s is not exist", recordId)
	}
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(recordId)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	var records []FinancialRecord
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var record FinancialRecord
		if err := json.Unmarshal(queryResult.Value, &record); err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}
