package fabric

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/hash"
	"time"
)

type ResultNotice struct {
	Tx     string `json:"tx"`
	Notice Notice `json:"notice"`
}

// Notice 公告信息
type Notice struct {
	MessageHash string `json:"message"`   //公告信息hash值
	Publisher   string `json:"publisher"` //公告发布人
	LastTime    string `json:"time"`      //上次更新时间
}

func CreateNotice(id, noticeHash, publisher string) error {
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
	contract := network.GetContract(noticeChaincode)

	_, err = contract.SubmitTransaction("CreateNotice", id, noticeHash, publisher)
	if err != nil {
		return fmt.Errorf("failed to submit transaction:%s", err.Error())
	}
	return nil
}

func UpdateNotice(id, noticeHash, publisher string) error {
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
	contract := network.GetContract(noticeChaincode)

	_, err = contract.SubmitTransaction("UpdateNotice", id, noticeHash, publisher)
	if err != nil {
		return fmt.Errorf("failed to submit transaction:%s", err.Error())
	}
	return nil
}

func GetNoticeHistory(id string) ([]ResultNotice, error) {
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
	contract := network.GetContract(noticeChaincode)

	result, err := contract.EvaluateTransaction("GetHistory", id)
	if err != nil {
		return nil, fmt.Errorf("failed to submit transaction:%s", err.Error())
	}
	var notices []ResultNotice
	if err := json.Unmarshal(result, &notices); err != nil {
		return nil, fmt.Errorf("failed to submit transaction:%s", err.Error())
	}
	return notices, nil
}

func Verify(id, noticeHash string) (bool, error) {

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
		return false, fmt.Errorf("failed to connect:%s", err.Error())
	}
	defer gw.Close()

	network := gw.GetNetwork(channel)
	contract := network.GetContract(noticeChaincode)

	result, err := contract.EvaluateTransaction("Verify", id, noticeHash)
	if err != nil {
		return false, fmt.Errorf("failed to submit transaction:%s", err.Error())
	}
	var b bool
	if err := json.Unmarshal(result, &b); err != nil {
		return false, fmt.Errorf("failed to submit transaction:%s", err.Error())
	}
	return b, nil
}
