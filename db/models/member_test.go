package models

import (
	"community-governance/db"
	"testing"
)

func TestLogin(t *testing.T) {
	err := db.InitDB()
	if err != nil {
		t.Error(err)
	}
	login, err := Login("110101199001010000", "123456")
	if err != nil {
		t.Error(err)
	}
	t.Log(login)
}
