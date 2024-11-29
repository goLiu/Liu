package fabric

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/hash"
	"time"
)

type Financial struct {
	MessageHash string `json:"message"`     //信息hash值
	Principal   string `json:"principal"`   //负责人ID
	UpdateDate  string `json:"update_date"` //上次修改时间
	State       string `json:"state"`       //款项状态
}
type FinancialRecord struct {
	Type       string `json:"type"`      //记录类型，0-收入 1-支出
	Amount     string `json:"amount"`    //记录金额
	Source     string `json:"source"`    //金额来源
	InfoHash   string `json:"info_hash"` //证明文件
	Explain    string `json:"explain"`   //说明
	RecorderID string `json:"record_id"` //记录人ID
}

const (
	foundStateActive = "active"
	foundStateClose  = "1"
)

type ChainFundDetail struct {
	BashHistory   []Financial       `json:"bash_history"`
	RecordHistory []FinancialRecord `json:"record_history"`
}

func CreateFinancial(id, finHash string) error {
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
	contract := network.GetContract(financialChaincode)

	_, err = contract.SubmitTransaction("CreateFinancial", id, finHash)
	if err != nil {
		return fmt.Errorf("failed to submit transaction:%s", err.Error())
	}
	return nil
}
func UpdateFinancial(id, finHash, state string) error {
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
	contract := network.GetContract(financialChaincode)

	_, err = contract.SubmitTransaction("UpdateFinancial", id, finHash)
	if err != nil {
		return fmt.Errorf("failed to submit UpdateFinancial transaction:%s", err.Error())
	}
	if state != foundStateActive {
		_, err := contract.SubmitTransaction("ExchangeState", id, foundStateClose)
		if err != nil {
			return fmt.Errorf("failed to submit ExchangeState transaction:%s", err.Error())
		}
	}
	return nil
}

func AddFinancialRecord(id, finType, source, explain, recorder, amount string) error {
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
	contract := network.GetContract(financialChaincode)

	_, err = contract.SubmitTransaction("AddRecord", id, finType, source, explain, recorder, amount)
	if err != nil {
		return fmt.Errorf("failed to submit transaction:%s", err.Error())
	}
	return nil
}
func ExchangeState(id, state string) error {
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
	contract := network.GetContract(financialChaincode)

	_, err = contract.SubmitTransaction("ExchangeState", id, state)
	if err != nil {
		return fmt.Errorf("failed to submit transaction:%s", err.Error())
	}
	return nil

}

func GetFinancialDetail(id string) (ChainFundDetail, error) {
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
		return ChainFundDetail{}, fmt.Errorf("failed to connect:%s", err.Error())
	}
	defer gw.Close()
	var detail ChainFundDetail

	network := gw.GetNetwork(channel)
	contract := network.GetContract(financialChaincode)

	result, err := contract.EvaluateTransaction("GetFinancialHistory", id)
	if err != nil {
		return ChainFundDetail{}, fmt.Errorf("failed to submit transaction:%s", err.Error())
	}
	if result != nil {
		if err := json.Unmarshal(result, &detail.BashHistory); err != nil {
			return ChainFundDetail{}, fmt.Errorf("failed to submit transaction:%s", err.Error())
		}
	}
	result, err = contract.EvaluateTransaction("GetFinancialRecordHistory", id)
	if err != nil {
		return ChainFundDetail{}, fmt.Errorf("failed to submit transaction:%s", err.Error())
	}
	if result != nil {
		if err := json.Unmarshal(result, &detail.RecordHistory); err != nil {
			return ChainFundDetail{}, fmt.Errorf("failed to submit transaction:%s", err.Error())
		}
	}
	return detail, nil
}
