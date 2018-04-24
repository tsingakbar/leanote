package tests

import (
	"github.com/tsingakbar/leanote/app/db"
	"testing"
	//	. "github.com/tsingakbar/leanote/app/lea"
	//	"github.com/tsingakbar/leanote/app/service"
	//	"gopkg.in/mgo.v2"
	//	"fmt"
)

func TestDBConnect(t *testing.T) {
	db.Init("mongodb://localhost:27017/leanote", "leanote")
}
