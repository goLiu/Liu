package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// FacilityContract 合约结构体
type FacilityContract struct {
	contractapi.Contract
}

type Facility struct {
	MessageHash string `json:"message"`     //信息hash值
	UpdateDate  string `json:"update_date"` //上次修改时间
	State       string `json:"state"`       //状态
}

// UsageRecord 使用记录结构体
type UsageRecord struct {
	User      string `json:"member"`    //借用人
	OperaTime string `json:"startTime"` //操作时间
	Operation string `json:"operation"` //操作
}

type Info struct {
	Facility Facility    `json:"facility"`
	Record   UsageRecord `json:"record"`
}

const (
	operationBorrow = "borrow"
	operationReturn = "return"
	stateAva        = "available"
	//stateRun  = "run"
	stateStop = "stop"
	recordFix = "record_"
)

// RegisterFacility 注册新的公共设施
func (f *FacilityContract) RegisterFacility(ctx contractapi.TransactionContextInterface, facilityID, messageHash string) error {
	//判断记录是否已经存在
	state, err := ctx.GetStub().GetState(facilityID)
	if err != nil {
		return fmt.Errorf("failed to get state:%s", err.Error())
	}
	if state != nil {
		return fmt.Errorf("%s is already exist", facilityID)
	}
	nowTime, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("failed to get tx timestamp:%s", err.Error())
	}
	//注册新公共设施
	facility := Facility{
		MessageHash: messageHash,
		UpdateDate:  nowTime.AsTime().Format("2006-01-02 15:04:05"),
		State:       stateAva,
	}
	facilityJSON, err := json.Marshal(facility)
	if err != nil {
		return fmt.Errorf("failed to marshal facility: %v", err)
	}
	return ctx.GetStub().PutState(facilityID, facilityJSON)
}

// GetFacility 获取Facility信息
func (f *FacilityContract) GetFacility(ctx contractapi.TransactionContextInterface, facilityID string) (Facility, error) {
	facilityJSON, err := ctx.GetStub().GetState(facilityID)
	if err != nil {
		return Facility{}, fmt.Errorf("failed to  get facility: %v", err.Error())
	}
	if facilityJSON == nil {
		return Facility{}, fmt.Errorf("%s is not exist", facilityID)
	}
	var facility Facility
	if err := json.Unmarshal(facilityJSON, &facility); err != nil {
		return Facility{}, fmt.Errorf("failed to un marshal:%s", err.Error())
	}
	return facility, nil
}

// RequestFacility 申请使用公共设施
func (f *FacilityContract) RequestFacility(ctx contractapi.TransactionContextInterface, facilityID, user string) error {
	//获取facility信息
	facility, err := f.GetFacility(ctx, facilityID)
	if err != nil {
		return err
	}
	//获取record信息
	id := recordFix + facilityID

	//判断用户状态是否可用
	if facility.State != stateAva {
		return fmt.Errorf("facility %s is currently %s", facilityID, facility.State)
	}
	//更新状态
	facility.State = stateStop
	updatedFacilityJSON, err := json.Marshal(facility)
	if err != nil {
		return fmt.Errorf("failed to marshal updated facility: %s", err.Error())
	}
	if err := ctx.GetStub().PutState(facilityID, updatedFacilityJSON); err != nil {
		return fmt.Errorf("failed to update facility state: %v", err)
	}
	stamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("failed to get time stamp:%s", err.Error())
	}

	// 记录使用信息
	usageRecord := UsageRecord{
		User:      user,
		OperaTime: stamp.AsTime().Format("2006-01-02 15:04:05"),
		Operation: operationBorrow,
	}
	recordJSON, err := json.Marshal(usageRecord)
	if err != nil {
		return fmt.Errorf("failed to marshal usage record: %v", err)
	}
	return ctx.GetStub().PutState(id, recordJSON)
}

// ReleaseFacility 释放公共设施
func (f *FacilityContract) ReleaseFacility(ctx contractapi.TransactionContextInterface, facilityID, user string) error {
	facility, err := f.GetFacility(ctx, facilityID)
	if err != nil {
		return err
	}
	// 更新设施状态为“可用”
	facility.State = stateAva
	updatedFacilityJSON, err := json.Marshal(facility)
	if err != nil {
		return fmt.Errorf("failed to marshal updated facility: %s", err.Error())
	}
	if err := ctx.GetStub().PutState(facilityID, updatedFacilityJSON); err != nil {
		return fmt.Errorf("failed to update facility state: %v", err)
	}
	//获取record信息
	id := recordFix + facilityID
	stamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("failed to get time stamp:%s", err.Error())
	}
	// 记录使用信息
	usageRecord := UsageRecord{
		User:      user,
		OperaTime: stamp.AsTime().Format("2006-01-02 15:04:05"),
		Operation: operationBorrow,
	}
	usageRecordJSON, err := json.Marshal(usageRecord)
	if err != nil {
		return fmt.Errorf("failed to marshal usage record: %v", err)
	}
	return ctx.GetStub().PutState(id, usageRecordJSON)
}

// GetFacilityUsageHistory 查询设施的使用记录
func (f *FacilityContract) GetFacilityUsageHistory(ctx contractapi.TransactionContextInterface, facilityID string) ([]UsageRecord, error) {
	id := recordFix + facilityID
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get history for %s: %v", facilityID, err)
	}
	defer resultsIterator.Close()

	var records []UsageRecord
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var record UsageRecord
		if err := json.Unmarshal(response.Value, &record); err == nil {
			records = append(records, record)
		}
	}
	return records, nil
}

func (f *FacilityContract) UpdateFacility(ctx contractapi.TransactionContextInterface, facilityID, messageHash, state string) error {
	facility, err := f.GetFacility(ctx, facilityID)
	if err != nil {
		return err
	}
	facility.State = state
	if state != "use" || facility.State == stateAva {
		facility.State = stateStop
	}
	facility.MessageHash = messageHash
	nowTime, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("failed to get time stamp:%s", err.Error())
	}
	facility.UpdateDate = nowTime.AsTime().Format("2006-01-02 15:04:05")
	updatedFacilityJSON, err := json.Marshal(facility)
	if err != nil {
		return fmt.Errorf("failed to marshal updated facility: %s", err.Error())
	}
	if err := ctx.GetStub().PutState(facilityID, updatedFacilityJSON); err != nil {
		return fmt.Errorf("failed to update facility state: %v", err)
	}
	return nil
}
