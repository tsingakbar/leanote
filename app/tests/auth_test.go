package tests

import (
	"github.com/tsingakbar/leanote/app/db"
	"testing"
	//	. "github.com/tsingakbar/leanote/app/lea"
	"github.com/tsingakbar/leanote/app/service"
	//	"gopkg.in/mgo.v2"
	//	"fmt"
)

func init() {
	db.Init("mongodb://localhost:27017/leanote", "leanote")
	service.InitService()
}

// 测试登录
func TestAuth(t *testing.T) {
	_, err := service.AuthS.Login("admin", "abc123")
	if err != nil {
		t.Error("Admin User Auth Error")
	}
}
