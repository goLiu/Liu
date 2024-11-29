package fabric

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/hash"
	"time"
)

type ChainVoteDetail struct {
	VoteNumber map[string]int `json:"vote_number"`
	Records    []VoteRecord   `json:"records"`
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

func CreatVote(id, base, ruleType, ruleValue, options string) error {
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
	contract := network.GetContract(voteChaincode)
	_, err = contract.SubmitTransaction("CreatVote", id, base, ruleType, ruleValue, options)
	if err != nil {
		return fmt.Errorf("failed to submit transaction:%s", err.Error())
	}
	return nil
}

func VoteJoin(id, voter, option string) (string, error) {
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
		return "", fmt.Errorf("failed to connect:%s", err.Error())
	}
	defer gw.Close()

	network := gw.GetNetwork(channel)
	contract := network.GetContract(voteChaincode)
	result, err := contract.SubmitTransaction("VoteJoin", id, voter, option)
	if err != nil {
		return "", fmt.Errorf("failed to submit transaction:%s", err.Error())
	}
	return string(result), nil
}

func GetVoteRecordDetail(id string) (ChainVoteDetail, error) {
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
		return ChainVoteDetail{}, fmt.Errorf("failed to connect:%s", err.Error())
	}
	defer gw.Close()

	network := gw.GetNetwork(channel)
	contract := network.GetContract(voteChaincode)
	result, err := contract.EvaluateTransaction("GetVoteRecordHistory", id)
	if err != nil {
		return ChainVoteDetail{}, fmt.Errorf("failed to submit transaction:%s", err.Error())
	}
	var records []VoteRecord
	if result == nil {
		if err := json.Unmarshal(result, &records); err != nil {
			return ChainVoteDetail{}, err
		}
	}
	//获取投票信息
	result, err = contract.EvaluateTransaction("GetVote", id)
	if err != nil {
		return ChainVoteDetail{}, fmt.Errorf("failed to submit transaction:%s", err.Error())
	}
	var vote Vote
	if err := json.Unmarshal(result, &vote); err != nil {
		return ChainVoteDetail{}, err
	}
	return ChainVoteDetail{VoteNumber: vote.Options, Records: records}, nil
}

func EndVote(id string) ([]string, error) {
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
	contract := network.GetContract(voteChaincode)
	result, err := contract.EvaluateTransaction("EndVote", id)
	if err != nil {
		return nil, fmt.Errorf("failed to submit transaction:%s", err.Error())
	}
	var options []string
	if err := json.Unmarshal(result, &options); err != nil {
		return nil, err
	}
	return options, nil
}
func CloseVote(id string) error {
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
	contract := network.GetContract(voteChaincode)
	_, err = contract.SubmitTransaction("CloseVote", id)
	if err != nil {
		return fmt.Errorf("failed to submit transaction:%s", err.Error())
	}
	return nil
}
