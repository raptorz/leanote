package service

import (
	"errors"
	"fmt"
	"github.com/leanote/leanote/app/db"
	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
	"strconv"
	"strings"
)

type AuthService struct {
}

func (this *AuthService) Login(emailOrUsername, pwd string) (info.User, error) {
	emailOrUsername = strings.Trim(emailOrUsername, " ")
	userInfo := userService.GetUserInfoByAny(emailOrUsername)
	if userInfo.UserId == "" || !ComparePwd(pwd, userInfo.Pwd) {
		return userInfo, errors.New("wrong username or password")
	}
	return userInfo, nil
}

func (this *AuthService) Register(email, pwd, fromUserId string) (bool, string) {
	if userService.IsExistsUser(email) {
		return false, "userHasBeenRegistered-" + email
	}
	passwd := GenPwd(pwd)
	if passwd == "" {
		return false, "GenerateHash error"
	}
	user := info.User{UserId: db.NewUUID(), Email: email, Username: email, Pwd: passwd}
	if fromUserId != "" && IsObjectId(fromUserId) {
		user.FromUserId = fromUserId
	}
	return this.register(user)
}

func (this *AuthService) register(user info.User) (bool, string) {
	if userService.AddUser(user) {
		userId := user.UserId
		notebook := info.Notebook{
			Seq:    -1,
			UserId: userId}
		title2Id := map[string]string{"life": db.NewUUID(), "study": db.NewUUID(), "work": db.NewUUID()}
		for title, objectId := range title2Id {
			notebook.Title = title
			notebook.NotebookId = objectId
			notebook.UserId = userId
			notebookService.AddNotebook(notebook)
		}

		registerSharedUserId := configService.GetGlobalStringConfig("registerSharedUserId")
		if registerSharedUserId != "" {
			registerSharedNotebooks := configService.GetGlobalArrMapConfig("registerSharedNotebooks")
			registerSharedNotes := configService.GetGlobalArrMapConfig("registerSharedNotes")
			registerCopyNoteIds := configService.GetGlobalArrayConfig("registerCopyNoteIds")

			for _, notebook := range registerSharedNotebooks {
				perm, _ := strconv.Atoi(notebook["perm"])
				shareService.AddShareNotebookToUserId(notebook["notebookId"], perm, registerSharedUserId, userId)
			}

			for _, note := range registerSharedNotes {
				perm, _ := strconv.Atoi(note["perm"])
				shareService.AddShareNoteToUserId(note["noteId"], perm, registerSharedUserId, userId)
			}

			for _, noteId := range registerCopyNoteIds {
				note := noteService.CopySharedNote(noteId, title2Id["life"], registerSharedUserId, userId)
				noteUpdate := map[string]interface{}{"IsBlog": false}
				noteService.UpdateNote(userId, note.NoteId, noteUpdate, -1)
			}
		}

		blogService.UpdateUserBlog(info.UserBlog{UserId: userId,
			Title:      user.Username + " 's Blog",
			SubTitle:   "Love Leanote!",
			AboutMe:    "Hello, I am (^_^)",
			CanComment: true,
		})
		blogService.AddOrUpdateSingle(userId, "", "About Me", "Hello, I am (^_^)")
	}

	return true, ""
}

func (this *AuthService) getUsername(thirdType, thirdUsername string) (username string) {
	username = thirdType + "-" + thirdUsername
	i := 1
	for {
		if !userService.IsExistsUserByUsername(username) {
			return
		}
		username = fmt.Sprintf("%v%v", username, i)
	}
}

func (this *AuthService) ThirdRegister(thirdType, thirdUserId, thirdUsername string) (exists bool, userInfo info.User) {
	userInfo = userService.GetUserInfoByThirdUserId(thirdUserId)
	if userInfo.UserId != "" {
		exists = true
		return
	}

	username := this.getUsername(thirdType, thirdUsername)
	userInfo = info.User{UserId: db.NewUUID(),
		Username:      username,
		ThirdUserId:   thirdUserId,
		ThirdUsername: thirdUsername,
	}
	_, _ = this.register(userInfo)
	return
}

func (this *AuthService) ActiveEmail(token string) (ok bool, msg, email string) {
	tokenInfo := info.Token{}
	if ok, msg, tokenInfo = tokenService.VerifyToken(token, info.TokenActiveEmail); ok {
		email = tokenInfo.Email
		userInfo := userService.GetUserInfoByEmail(email)
		if userInfo.UserId == "" {
			ok = false
			msg = "user not exists"
			return
		}

		_, err := db.DB.Exec("UPDATE users SET verified = true WHERE id = $1", userInfo.UserId)
		if err != nil {
			ok = false
			msg = "database error"
		}
		return
	}

	ok = false
	msg = "token expired"
	return
}

func (this *AuthService) UpdateEmail(token string) (ok bool, msg, email string) {
	tokenInfo := info.Token{}
	if ok, msg, tokenInfo = tokenService.VerifyToken(token, info.TokenUpdateEmail); ok {
		email = strings.ToLower(tokenInfo.Email)
		if userService.IsExistsUser(email) {
			ok = false
			msg = "email already registered"
			return
		}

		_, err := db.DB.Exec("UPDATE users SET email = $1, verified = true WHERE id = $2", email, tokenInfo.UserId)
		if err != nil {
			ok = false
			msg = "database error"
		}
		return
	}

	ok = false
	msg = "token expired"
	return
}
