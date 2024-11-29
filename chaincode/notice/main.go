package main

import (
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	"log"
)

func main() {
	noticeContract, err := contractapi.NewChaincode(&NoticeContract{})
	if err != nil {
		log.Printf("error create notice contract:%s", err.Error())
	}
	if err := noticeContract.Start(); err != nil {
		log.Printf("failed start chaincode:%s", err.Error())
	}
}
