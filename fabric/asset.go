package fabric

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/hash"
	"time"
)

type Asset struct {
	AssetHash  string `json:"assetHash"` // 资产HASH值
	Owner      string `json:"owner"`     //资产拥有者
	Recorder   string `json:"recorder"`
	RecordDate string `json:"recordDate"`
}

func CreateAsset(assetID, asserHash, owner, recorder string) error {
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	identity := newIdentity()
	sign := newSign()

	gw, err := client.Connect(
		identity,
		client.WithSign(sign),
		client.WithHash(hash.SHA256),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		return fmt.Errorf("failed to connect:%s", err.Error())
	}
	defer gw.Close()

	network := gw.GetNetwork(channel)
	contract := network.GetContract(assetChaincode)
	_, err = contract.SubmitTransaction("CreateAsset", assetID, asserHash, owner, recorder)
	if err != nil {
		return fmt.Errorf("failed to submit transaction:%s", err.Error())
	}
	return nil
}

// func GetAsset() {
//
// }
func ExchangeOwner(assetID string, newOwner string) error {
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	identity := newIdentity()
	sign := newSign()

	gw, err := client.Connect(
		identity,
		client.WithSign(sign),
		client.WithHash(hash.SHA256),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		return fmt.Errorf("failed to connect:%s", err.Error())
	}
	defer gw.Close()

	network := gw.GetNetwork(channel)
	contract := network.GetContract(assetChaincode)
	_, err = contract.SubmitTransaction("ExchangeOwner", assetID, newOwner)
	if err != nil {
		return fmt.Errorf("failed to submit transaction:%s", err.Error())
	}
	return nil
}
func GetAssetHistory(assetID string) ([]Asset, error) {
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	identity := newIdentity()
	sign := newSign()

	gw, err := client.Connect(
		identity,
		client.WithSign(sign),
		client.WithHash(hash.SHA256),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect:%s", err.Error())
	}
	defer gw.Close()

	network := gw.GetNetwork(channel)
	contract := network.GetContract(assetChaincode)

	result, err := contract.EvaluateTransaction("GetAssetHistory", assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction:%s", err.Error())
	}
	var assetList []Asset
	err = json.Unmarshal(result, &assetList)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal result:%s", err.Error())
	}

	return assetList, nil
}
func UpdateAsset(assetID, asserHash, owner, recorder string) error {
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	identity := newIdentity()
	sign := newSign()

	gw, err := client.Connect(
		identity,
		client.WithSign(sign),
		client.WithHash(hash.SHA256),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		return fmt.Errorf("failed to connect:%s", err.Error())
	}
	defer gw.Close()

	network := gw.GetNetwork(channel)
	contract := network.GetContract(assetChaincode)
	_, err = contract.SubmitTransaction("UpdateAsset", assetID, asserHash, owner, recorder)
	if err != nil {
		return err
	}
	return nil
}
