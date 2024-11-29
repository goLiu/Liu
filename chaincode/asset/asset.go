package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// AssetContract  资产管理合约
type AssetContract struct {
	contractapi.Contract
}
type Asset struct {
	AssetHash  string `json:"assetHash"` // 资产HASH值
	Owner      string `json:"owner"`     //资产拥有者
	Recorder   string `json:"recorder"`
	RecordDate string `json:"recordDate"`
}

const (
	stateInit        = "active"  // 资产状态 - 初始化
	stateMaintenance = "repair"  // 资产状态 - 维修中
	stateExpired     = "disposa" // 资产状态 - 失效
)

// CreateAsset 创建资产
func (a *AssetContract) CreateAsset(ctx contractapi.TransactionContextInterface, assetID, asserHash, owner, recorder string) error {
	state, err := ctx.GetStub().GetState(assetID)
	if err != nil {
		return fmt.Errorf("failed get state:%s", err.Error())
	}
	if state != nil {
		return fmt.Errorf("%s already exist", assetID)
	}
	nowTime, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("failed to get tx timestamp:%s", err.Error())
	}
	asset := Asset{
		AssetHash:  asserHash,
		Owner:      owner,
		Recorder:   recorder,
		RecordDate: nowTime.AsTime().Format("2006-01-02 15:04:05"),
	}
	assetBytes, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf("failed marshal asset:%s", err.Error())
	}
	return ctx.GetStub().PutState(assetID, assetBytes)
}

// GetAsset  根据id查询指定资产信息
func (a *AssetContract) GetAsset(ctx contractapi.TransactionContextInterface, assetID string) (Asset, error) {
	if assetID == "" {
		return Asset{}, fmt.Errorf("id = ''")
	}
	state, err := ctx.GetStub().GetState(assetID)
	if err != nil {
		return Asset{}, fmt.Errorf("failed to get state by %s :%s", assetID, err.Error())
	}
	if state == nil {
		return Asset{}, fmt.Errorf("%s not exist", assetID)
	}
	var asset Asset
	if err := json.Unmarshal(state, &asset); err != nil {
		return Asset{}, fmt.Errorf("failed to unmarshar state:%s", err.Error())
	}
	return asset, nil
}

// ExchangeOwner 更改资产owner
func (a *AssetContract) ExchangeOwner(ctx contractapi.TransactionContextInterface, assetID string, newOwner string) error {
	//判断assetID是否存在
	asset, err := a.GetAsset(ctx, assetID)
	if err != nil {
		return err
	}
	asset.Owner = newOwner
	assetBytes, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf("failed to marashl:%s", err.Error())
	}
	return ctx.GetStub().PutState(assetID, assetBytes)
}

// GetAssetHistory 查询资产变更记录
func (a *AssetContract) GetAssetHistory(ctx contractapi.TransactionContextInterface, assetID string) ([]Asset, error) {
	state, err := ctx.GetStub().GetHistoryForKey(assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get history record for %s: %v", assetID, err)
	}
	defer state.Close()

	var assets []Asset
	for state.HasNext() {
		next, err := state.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next history record: %v", err)
		}

		var asset Asset
		if err := json.Unmarshal(next.Value, &asset); err != nil {
			return nil, fmt.Errorf("failed to unmarshal asset data: %v", err)
		}

		assets = append(assets, asset)
	}

	return assets, nil
}

// UpdateAsset 更改资产信息
func (a *AssetContract) UpdateAsset(ctx contractapi.TransactionContextInterface, assetID, asserHash, owner, recorder string) error {
	asset, err := a.GetAsset(ctx, assetID)
	if err != nil {
		return err
	}
	nowTime, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("failed to get tx timestamp:%s", err.Error())
	}
	asset.RecordDate = nowTime.AsTime().Format("2006-01-02 15:04:05")
	asset.AssetHash = asserHash
	asset.Owner = owner
	asset.Recorder = recorder
	assetBytes, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf("failed marshal asset:%s", err.Error())
	}
	return ctx.GetStub().PutState(assetID, assetBytes)
}
