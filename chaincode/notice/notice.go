package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// NoticeContract 公告管理合约
type NoticeContract struct {
	contractapi.Contract
}

// Notice 公告信息
type Notice struct {
	MessageHash string `json:"message"`   //公告信息hash值
	Publisher   string `json:"publisher"` //公告发布人
	LastTime    string `json:"time"`      //上次更新时间
}

// CreateNotice 上传公告信息
func (n *NoticeContract) CreateNotice(ctx contractapi.TransactionContextInterface, id, hash, publisher string) error {
	//检查id是否已经存在
	state, err := ctx.GetStub().GetState(id)
	if err != nil {
		return err
	}
	if state != nil {
		return fmt.Errorf("%s already exist", id)
	}
	nowStamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return err
	}
	notice := Notice{
		MessageHash: hash,
		Publisher:   publisher,
		LastTime:    nowStamp.AsTime().Format("2006-01-02 15:04:05"),
	}
	infoBytes, err := json.Marshal(notice)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(id, infoBytes)
}

// UpdateNotice  更新公告信息
func (n *NoticeContract) UpdateNotice(ctx contractapi.TransactionContextInterface, id, hash, publisher string) error {
	state, err := ctx.GetStub().GetState(id)
	if err != nil {
		return err
	}
	if state == nil {
		return fmt.Errorf("%s not exist", id)
	}
	var notice Notice
	if err := json.Unmarshal(state, &notice); err != nil {
		return err
	}
	notice.MessageHash = hash
	notice.Publisher = publisher
	nowStamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return err
	}
	notice.LastTime = nowStamp.AsTime().Format("2006-01-02 15:04:05")
	infoBytes, err := json.Marshal(notice)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(id, infoBytes)
}

// GetNotice 查询指定公告信息
func (n *NoticeContract) GetNotice(ctx contractapi.TransactionContextInterface, id string) (Notice, error) {
	state, err := ctx.GetStub().GetState(id)
	if err != nil {
		return Notice{}, err
	}
	if state == nil {
		return Notice{}, fmt.Errorf("%s not exist", id)
	}
	var notice Notice
	if err := json.Unmarshal(state, &notice); err != nil {
		return Notice{}, err
	}
	return notice, nil
}

// Verify  验证公告信息是否与链下信息一致
func (n *NoticeContract) Verify(ctx contractapi.TransactionContextInterface, id, hash string) (bool, error) {
	state, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, err
	}
	if state == nil {
		return false, fmt.Errorf("%s not exist", id)
	}

	var notice Notice
	if err := json.Unmarshal(state, &notice); err != nil {
		return false, err
	}
	return notice.MessageHash == hash, nil
}

type ResultNotice struct {
	Tx     string `json:"tx"`
	Notice Notice `json:"notice"`
}

// GetHistory 查看指定公告的记录变更历史
func (n *NoticeContract) GetHistory(ctx contractapi.TransactionContextInterface, id string) ([]ResultNotice, error) {
	var results []ResultNotice
	iterator, err := ctx.GetStub().GetHistoryForKey(id)
	if err != nil {
		return nil, err
	}
	defer iterator.Close()
	for iterator.HasNext() {
		var info ResultNotice
		nextInfo, err := iterator.Next()
		if err != nil {
			return nil, err
		}
		info.Tx = nextInfo.TxId
		if err := json.Unmarshal(nextInfo.Value, &info.Notice); err != nil {
			return nil, fmt.Errorf("failed to convert notice %e", err)
		}
		results = append(results, info)
	}
	return results, nil
}
