package service

import (
	"github.com/leanote/leanote/app/db"
	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
)

// 找回密码
// 修改密码
var overHours = 2.0 // 小时后过期

type PwdService struct {
}

// 1. 找回密码, 通过email找用户,
// 用户存在, 生成code
func (this *PwdService) FindPwd(email string) (ok bool, msg string) {
	ok = false
	userId := userService.GetUserId(email)
	if userId == "" {
		msg = "用户不存在"
		return
	}

	token := tokenService.NewToken(userId, email, info.TokenPwd)
	if token == "" {
		return false, "db error"
	}

	// 发送邮件
	// TODO: 实现邮件发送功能
	// ok, msg = emailService.FindPwdSendEmail(token, email)
	ok, msg = false, "邮件功能暂不可用"
	return
}

// 重置密码时
// 修改密码
// 先验证
func (this *PwdService) UpdatePwd(token, pwd string) (bool, string) {
	var tokenInfo info.Token
	var ok bool
	var msg string

	// 先验证
	if ok, msg, tokenInfo = tokenService.VerifyToken(token, info.TokenPwd); !ok {
		return ok, msg
	}

	passwd := GenPwd(pwd)
	if passwd == "" {
		return false, "GenerateHash error"
	}

	// 修改密码之
	query := "UPDATE users SET pwd = $1 WHERE id = $2"
	_, err := db.DB.Exec(query, passwd, tokenInfo.UserId)
	ok = err == nil

	// 删除token
	tokenService.DeleteToken(tokenInfo.UserId, info.TokenPwd)

	return ok, ""
}
