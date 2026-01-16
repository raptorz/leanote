package service

import (
	"database/sql"
	"github.com/leanote/leanote/app/db"
	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
	"time"
)

// Session存储到mongodb中
type SessionService struct {
}

func (this *SessionService) Update(sessionId, key string, value interface{}) bool {
	// TODO: 实现PostgreSQL版本
	// 由于session值可能是任意类型，这里简化处理
	query := "UPDATE sessions SET updated_time = $1 WHERE session_id = $2"
	_, err := db.DB.Exec(query, time.Now(), sessionId)
	return err == nil
}

// 注销时清空session
func (this *SessionService) Clear(sessionId string) bool {
	query := "DELETE FROM sessions WHERE session_id = $1"
	_, err := db.DB.Exec(query, sessionId)
	return err == nil
}
func (this *SessionService) Get(sessionId string) info.Session {
	session := info.Session{}
	query := "SELECT id, session_id, created_time, updated_time FROM sessions WHERE session_id = $1"
	err := db.DB.QueryRow(query, sessionId).Scan(&session.Id, &session.SessionId, &session.CreatedTime, &session.UpdatedTime)

	// 如果没有session, 那么插入一条之
	if err == sql.ErrNoRows {
		session.Id = db.NewUUID()
		session.SessionId = sessionId
		session.CreatedTime = time.Now()
		session.UpdatedTime = session.CreatedTime
		insertQuery := "INSERT INTO sessions (id, session_id, created_time, updated_time) VALUES ($1, $2, $3, $4)"
		_, err := db.DB.Exec(insertQuery, session.Id, session.SessionId, session.CreatedTime, session.UpdatedTime)
		if err != nil {
			Log(err.Error())
		}
	}

	return session
}

//------------------
// 错误次数处理

// 登录错误时间是否已超过了
func (this *SessionService) LoginTimesIsOver(sessionId string) bool {
	session := this.Get(sessionId)
	return session.LoginTimes > 5
}

// 登录成功后清空错误次数
func (this *SessionService) ClearLoginTimes(sessionId string) bool {
	return this.Update(sessionId, "LoginTimes", 0)
}

// 增加错误次数
func (this *SessionService) IncrLoginTimes(sessionId string) bool {
	session := this.Get(sessionId)
	return this.Update(sessionId, "LoginTimes", session.LoginTimes+1)
}

// ----------
// 验证码
func (this *SessionService) GetCaptcha(sessionId string) string {
	session := this.Get(sessionId)
	return session.Captcha
}
func (this *SessionService) SetCaptcha(sessionId, captcha string) bool {
	this.Get(sessionId)
	// Log(sessionId)
	// Log(captcha)
	ok := this.Update(sessionId, "Captcha", captcha)
	// Log(ok)
	return ok
}

// -----------
// API
func (this *SessionService) GetUserId(sessionId string) string {
	session := this.Get(sessionId)
	// 更新updateTime, 避免过期
	query := "UPDATE sessions SET updated_time = $1 WHERE session_id = $2"
	_, _ = db.DB.Exec(query, time.Now(), sessionId)
	return session.UserId
}

// 登录成功后设置userId
func (this *SessionService) SetUserId(sessionId, userId string) bool {
	this.Get(sessionId)
	ok := this.Update(sessionId, "UserId", userId)
	return ok
}
