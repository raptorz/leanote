package tests

import (
	"testing"
	//	. "github.com/leanote/leanote/app/lea"
	"github.com/leanote/leanote/app/service"
	//	"gopkg.in/mgo.v2"
	//	"fmt"
)

// 测试登录
func TestSendMail(t *testing.T) {
	requireMongoIntegration(t)
	ok, err := service.EmailS.SendEmail("life@leanote.com", "你好", "你好吗")
	t.Log(ok)
	t.Log(err)
}
