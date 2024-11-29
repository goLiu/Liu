package models

import (
	"community-governance/db"
	"testing"
)

func TestGetVoteRuleById(t *testing.T) {
	db.InitDB()
	ruleId, err := GetRuleIdByName("最低投票要求")
	if err != nil {
		t.Errorf("GetRuleIdByName failed: %v", err)
		return
	}

	t.Logf("GetRuleIdByName success: %s", ruleId)
}
