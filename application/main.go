package main

import (
	"community-governance/application/router"
	"community-governance/db"
	"log"
)

func main() {
	//初始化db
	err := db.InitDB()
	if err != nil {
		log.Fatalf("failed to init db:%s", err.Error())
	}
	r := router.SetupRouter()

	// 启动服务器
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
