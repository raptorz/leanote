package service

import (
	"github.com/leanote/leanote/app/db"
	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
	"time"
)

// token
// 找回密码
// 修改密码

type TokenService struct {
}

// 生成token
func (this *TokenService) NewToken(userId string, email string, tokenType int) string {
	token := info.Token{UserId: userId, Token: NewGuidWith(email), CreatedTime: time.Now(), Email: email, Type: tokenType}

	// PostgreSQL UPSERT实现
	query := `
		INSERT INTO tokens (user_id, token, created_time, email, type) 
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id) 
		DO UPDATE SET token = $2, created_time = $3, email = $4, type = $5
	`
	_, err := db.DB.Exec(query, token.UserId, token.Token, token.CreatedTime, token.Email, token.Type)
	if err == nil {
		return token.Token
	}

	return ""
}

// 删除token
func (this *TokenService) DeleteToken(userId string, tokenType int) bool {
	query := `DELETE FROM tokens WHERE user_id = $1 AND type = $2`
	_, err := db.DB.Exec(query, userId, tokenType)
	return err == nil
}

func (this *TokenService) GetOverHours(tokenType int) float64 {
	if tokenType == info.TokenPwd {
		return info.PwdOverHours
	} else if tokenType == info.TokenUpdateEmail {
		return info.UpdateEmailOverHours
	} else {
		return info.ActiveEmailOverHours
	}
}

// 验证token, 是否存在, 过时?
func (this *TokenService) VerifyToken(token string, tokenType int) (ok bool, msg string, tokenInfo info.Token) {
	overHours := this.GetOverHours(tokenType)

	ok = false
	if token == "" {
		msg = "不存在"
		return
	}

	// PostgreSQL查询
	query := `SELECT user_id, token, created_time, email, type FROM tokens WHERE token = $1`
	row := db.DB.QueryRow(query, token)

	err := row.Scan(&tokenInfo.UserId, &tokenInfo.Token, &tokenInfo.CreatedTime, &tokenInfo.Email, &tokenInfo.Type)
	if err != nil {
		msg = "不存在"
		return
	}

	if tokenInfo.UserId == "" {
		msg = "不存在"
		return
	}

	// 验证是否过时
	now := time.Now()
	duration := now.Sub(tokenInfo.CreatedTime)

	if duration.Hours() > overHours {
		msg = "过期"
		return
	}

	ok = true
	return
}
