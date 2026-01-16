package service

import (
	"database/sql"
	"github.com/leanote/leanote/app/db"
	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
	"github.com/lib/pq"
	"strings"
	"time"
)

type UserService struct {
}

func (this *UserService) IncrUsn(userId string) int {
	var usn int
	err := db.DB.QueryRow("SELECT usn FROM users WHERE id = $1", userId).Scan(&usn)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0
		}
		return 0
	}
	usn += 1
	_, err = db.DB.Exec("UPDATE users SET usn = $1 WHERE id = $2", usn, userId)
	if err != nil {
		return 0
	}
	return usn
}

func (this *UserService) GetUsn(userId string) int {
	var usn int
	err := db.DB.QueryRow("SELECT usn FROM users WHERE id = $1", userId).Scan(&usn)
	if err != nil {
		return 0
	}
	return usn
}

func (this *UserService) AddUser(user info.User) bool {
	if user.UserId == "" {
		user.UserId = db.NewUUID()
	}
	user.CreatedTime = time.Now()

	if user.Email != "" {
		user.Email = strings.ToLower(user.Email)

		go func() {
			emailService.RegisterSendActiveEmail(user, user.Email)
		}()
	}

	query := `INSERT INTO users (
		id, email, verified, username, username_raw, pwd, created_time,
		logo, theme, notebook_width, note_list_width, md_editor_width, left_is_min,
		third_user_id, third_username, third_type, image_num, image_size,
		attach_num, attach_size, from_user_id, account_type, account_start_time,
		account_end_time, max_image_num, max_image_size, max_attach_num,
		max_attach_size, max_per_attach_size, usn, full_sync_before, is_deleted
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
		$16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29,
		$30, $31
	)`

	_, err := db.DB.Exec(query,
		user.UserId,
		user.Email,
		user.Verified,
		user.Username,
		user.UsernameRaw,
		user.Pwd,
		user.CreatedTime,
		user.Logo,
		user.Theme,
		user.NotebookWidth,
		user.NoteListWidth,
		user.MdEditorWidth,
		user.LeftIsMin,
		user.ThirdUserId,
		user.ThirdUsername,
		user.ThirdType,
		user.ImageNum,
		user.ImageSize,
		user.AttachNum,
		user.AttachSize,
		user.FromUserId,
		user.AccountType,
		user.AccountStartTime,
		user.AccountEndTime,
		user.MaxImageNum,
		user.MaxImageSize,
		user.MaxAttachNum,
		user.MaxAttachSize,
		user.MaxPerAttachSize,
		user.Usn,
		user.FullSyncBefore,
		user.IsDeleted,
	)

	return err == nil
}

func (this *UserService) GetUserId(email string) string {
	email = strings.ToLower(email)
	var userId string
	err := db.DB.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&userId)
	if err != nil {
		return ""
	}
	return userId
}

func (this *UserService) GetUsername(userId string) string {
	var username string
	err := db.DB.QueryRow("SELECT username FROM users WHERE id = $1", userId).Scan(&username)
	if err != nil {
		return ""
	}
	return username
}

func (this *UserService) IsExistsUser(email string) bool {
	return this.GetUserId(email) != ""
}

func (this *UserService) IsExistsUserByUsername(username string) bool {
	var count int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", username).Scan(&count)
	if err != nil {
		return false
	}
	return count >= 1
}

func (this *UserService) GetUserInfoByAny(idEmailUsername string) info.User {
	if IsObjectId(idEmailUsername) {
		return this.GetUserInfo(idEmailUsername)
	}

	if strings.Contains(idEmailUsername, "@") {
		return this.GetUserInfoByEmail(idEmailUsername)
	}

	return this.GetUserInfoByUsername(idEmailUsername)
}

func (this *UserService) GetUserInfo(userId string) info.User {
	var user info.User
	query := `SELECT id, email, verified, username, username_raw, pwd, created_time,
		logo, theme, notebook_width, note_list_width, md_editor_width, left_is_min,
		third_user_id, third_username, third_type, image_num, image_size,
		attach_num, attach_size, from_user_id, account_type, account_start_time,
		account_end_time, max_image_num, max_image_size, max_attach_num,
		max_attach_size, max_per_attach_size, usn, full_sync_before, is_deleted
		FROM users WHERE id = $1`

	err := db.DB.QueryRow(query, userId).Scan(
		&user.UserId,
		&user.Email,
		&user.Verified,
		&user.Username,
		&user.UsernameRaw,
		&user.Pwd,
		&user.CreatedTime,
		&user.Logo,
		&user.Theme,
		&user.NotebookWidth,
		&user.NoteListWidth,
		&user.MdEditorWidth,
		&user.LeftIsMin,
		&user.ThirdUserId,
		&user.ThirdUsername,
		&user.ThirdType,
		&user.ImageNum,
		&user.ImageSize,
		&user.AttachNum,
		&user.AttachSize,
		&user.FromUserId,
		&user.AccountType,
		&user.AccountStartTime,
		&user.AccountEndTime,
		&user.MaxImageNum,
		&user.MaxImageSize,
		&user.MaxAttachNum,
		&user.MaxAttachSize,
		&user.MaxPerAttachSize,
		&user.Usn,
		&user.FullSyncBefore,
		&user.IsDeleted,
	)

	if err != nil {
		return info.User{}
	}

	return user
}

func (this *UserService) GetUserInfoByEmail(email string) info.User {
	email = strings.ToLower(email)
	var user info.User
	query := `SELECT id, email, verified, username, username_raw, pwd, created_time,
		logo, theme, notebook_width, note_list_width, md_editor_width, left_is_min,
		third_user_id, third_username, third_type, image_num, image_size,
		attach_num, attach_size, from_user_id, account_type, account_start_time,
		account_end_time, max_image_num, max_image_size, max_attach_num,
		max_attach_size, max_per_attach_size, usn, full_sync_before, is_deleted
		FROM users WHERE email = $1`

	err := db.DB.QueryRow(query, email).Scan(
		&user.UserId,
		&user.Email,
		&user.Verified,
		&user.Username,
		&user.UsernameRaw,
		&user.Pwd,
		&user.CreatedTime,
		&user.Logo,
		&user.Theme,
		&user.NotebookWidth,
		&user.NoteListWidth,
		&user.MdEditorWidth,
		&user.LeftIsMin,
		&user.ThirdUserId,
		&user.ThirdUsername,
		&user.ThirdType,
		&user.ImageNum,
		&user.ImageSize,
		&user.AttachNum,
		&user.AttachSize,
		&user.FromUserId,
		&user.AccountType,
		&user.AccountStartTime,
		&user.AccountEndTime,
		&user.MaxImageNum,
		&user.MaxImageSize,
		&user.MaxAttachNum,
		&user.MaxAttachSize,
		&user.MaxPerAttachSize,
		&user.Usn,
		&user.FullSyncBefore,
		&user.IsDeleted,
	)

	if err != nil {
		return info.User{}
	}

	return user
}

func (this *UserService) GetUserInfoByUsername(username string) info.User {
	var user info.User
	query := `SELECT id, email, verified, username, username_raw, pwd, created_time,
		logo, theme, notebook_width, note_list_width, md_editor_width, left_is_min,
		third_user_id, third_username, third_type, image_num, image_size,
		attach_num, attach_size, from_user_id, account_type, account_start_time,
		account_end_time, max_image_num, max_image_size, max_attach_num,
		max_attach_size, max_per_attach_size, usn, full_sync_before, is_deleted
		FROM users WHERE username = $1`

	err := db.DB.QueryRow(query, username).Scan(
		&user.UserId,
		&user.Email,
		&user.Verified,
		&user.Username,
		&user.UsernameRaw,
		&user.Pwd,
		&user.CreatedTime,
		&user.Logo,
		&user.Theme,
		&user.NotebookWidth,
		&user.NoteListWidth,
		&user.MdEditorWidth,
		&user.LeftIsMin,
		&user.ThirdUserId,
		&user.ThirdUsername,
		&user.ThirdType,
		&user.ImageNum,
		&user.ImageSize,
		&user.AttachNum,
		&user.AttachSize,
		&user.FromUserId,
		&user.AccountType,
		&user.AccountStartTime,
		&user.AccountEndTime,
		&user.MaxImageNum,
		&user.MaxImageSize,
		&user.MaxAttachNum,
		&user.MaxAttachSize,
		&user.MaxPerAttachSize,
		&user.Usn,
		&user.FullSyncBefore,
		&user.IsDeleted,
	)

	if err != nil {
		return info.User{}
	}

	return user
}

func (this *UserService) UpdateUser(user info.User) bool {
	query := `UPDATE users SET
		email = $2, verified = $3, username = $4, username_raw = $5, pwd = $6,
		logo = $7, theme = $8, notebook_width = $9, note_list_width = $10,
		md_editor_width = $11, left_is_min = $12, third_user_id = $13,
		third_username = $14, third_type = $15, image_num = $16, image_size = $17,
		attach_num = $18, attach_size = $19, from_user_id = $20, account_type = $21,
		account_start_time = $22, account_end_time = $23, max_image_num = $24,
		max_image_size = $25, max_attach_num = $26, max_attach_size = $27,
		max_per_attach_size = $28, usn = $29, full_sync_before = $30
		WHERE id = $1`

	_, err := db.DB.Exec(query,
		user.UserId,
		user.Email,
		user.Verified,
		user.Username,
		user.UsernameRaw,
		user.Pwd,
		user.Logo,
		user.Theme,
		user.NotebookWidth,
		user.NoteListWidth,
		user.MdEditorWidth,
		user.LeftIsMin,
		user.ThirdUserId,
		user.ThirdUsername,
		user.ThirdType,
		user.ImageNum,
		user.ImageSize,
		user.AttachNum,
		user.AttachSize,
		user.FromUserId,
		user.AccountType,
		user.AccountStartTime,
		user.AccountEndTime,
		user.MaxImageNum,
		user.MaxImageSize,
		user.MaxAttachNum,
		user.MaxAttachSize,
		user.MaxPerAttachSize,
		user.Usn,
		user.FullSyncBefore,
	)

	return err == nil
}

func (this *UserService) GetUserInfoByThirdUserId(thirdUserId string) info.User {
	var user info.User
	query := `SELECT id, email, verified, username, username_raw, pwd, created_time,
		logo, theme, notebook_width, note_list_width, md_editor_width, left_is_min,
		third_user_id, third_username, third_type, image_num, image_size,
		attach_num, attach_size, from_user_id, account_type, account_start_time,
		account_end_time, max_image_num, max_image_size, max_attach_num,
		max_attach_size, max_per_attach_size, usn, full_sync_before, is_deleted
		FROM users WHERE third_user_id = $1`

	err := db.DB.QueryRow(query, thirdUserId).Scan(
		&user.UserId,
		&user.Email,
		&user.Verified,
		&user.Username,
		&user.UsernameRaw,
		&user.Pwd,
		&user.CreatedTime,
		&user.Logo,
		&user.Theme,
		&user.NotebookWidth,
		&user.NoteListWidth,
		&user.MdEditorWidth,
		&user.LeftIsMin,
		&user.ThirdUserId,
		&user.ThirdUsername,
		&user.ThirdType,
		&user.ImageNum,
		&user.ImageSize,
		&user.AttachNum,
		&user.AttachSize,
		&user.FromUserId,
		&user.AccountType,
		&user.AccountStartTime,
		&user.AccountEndTime,
		&user.MaxImageNum,
		&user.MaxImageSize,
		&user.MaxAttachNum,
		&user.MaxAttachSize,
		&user.MaxPerAttachSize,
		&user.Usn,
		&user.FullSyncBefore,
		&user.IsDeleted,
	)

	if err != nil {
		return info.User{}
	}

	return user
}

func (this *UserService) UpdateUsername(userId, username string) (bool, string) {
	if userId == "" || username == "" || username == "admin" {
		return false, "usernameIsExisted"
	}
	usernameRaw := username
	username = strings.ToLower(username)

	var exists int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1 AND id != $2", username, userId).Scan(&exists)
	if err != nil {
		return false, "databaseError"
	}
	if exists > 0 {
		return false, "usernameIsExisted"
	}

	_, err = db.DB.Exec("UPDATE users SET username = $1, username_raw = $2 WHERE id = $3", username, usernameRaw, userId)
	return err == nil, ""
}

func (this *UserService) UpdateAvatar(userId, avatarPath string) bool {
	_, err := db.DB.Exec("UPDATE users SET logo = $1 WHERE id = $2", avatarPath, userId)
	return err == nil
}

func (this *UserService) UpdatePwd(userId, oldPwd, pwd string) (bool, string) {
	userInfo := this.GetUserInfo(userId)
	if !ComparePwd(oldPwd, userInfo.Pwd) {
		return false, "oldPasswordError"
	}

	passwd := GenPwd(pwd)
	if passwd == "" {
		return false, "GenerateHash error"
	}

	_, err := db.DB.Exec("UPDATE users SET pwd = $1 WHERE id = $2", passwd, userId)
	return err == nil, ""
}

func (this *UserService) ResetPwd(adminUserId, userId, pwd string) (ok bool, msg string) {
	if configService.GetAdminUserId() != adminUserId {
		return false, "notAdmin"
	}

	passwd := GenPwd(pwd)
	if passwd == "" {
		return false, "GenerateHash error"
	}
	_, err := db.DB.Exec("UPDATE users SET pwd = $1 WHERE id = $2", passwd, userId)
	return err == nil, ""
}

func (this *UserService) UpdateTheme(userId, theme string) bool {
	_, err := db.DB.Exec("UPDATE users SET theme = $1 WHERE id = $2", theme, userId)
	return err == nil
}

func (this *UserService) UpdateAccount(userId, accountType string, accountStartTime, accountEndTime time.Time,
	maxImageNum, maxImageSize, maxAttachNum, maxAttachSize, maxPerAttachSize int) bool {
	_, err := db.DB.Exec(`
		UPDATE users SET 
			account_type = $2,
			account_start_time = $3,
			account_end_time = $4,
			max_image_num = $5,
			max_image_size = $6,
			max_attach_num = $7,
			max_attach_size = $8,
			max_per_attach_size = $9
		WHERE id = $1`,
		userId, accountType, accountStartTime, accountEndTime,
		maxImageNum, maxImageSize, maxAttachNum, maxAttachSize, maxPerAttachSize)
	return err == nil
}

func (this *UserService) UpdateColumnWidth(userId string, notebookWidth, noteListWidth, mdEditorWidth int) bool {
	_, err := db.DB.Exec("UPDATE users SET notebook_width = $1, note_list_width = $2, md_editor_width = $3 WHERE id = $4",
		notebookWidth, noteListWidth, mdEditorWidth, userId)
	return err == nil
}

func (this *UserService) UpdateLeftIsMin(userId string, leftIsMin bool) bool {
	_, err := db.DB.Exec("UPDATE users SET left_is_min = $1 WHERE id = $2", leftIsMin, userId)
	return err == nil
}

func (this *UserService) MapUserInfoAndBlogInfosByUserIds(userIds []string) map[string]info.UserAndBlog {
	return make(map[string]info.UserAndBlog)
}

func (this *UserService) MapUserAndBlogByUserIds(userIds []string) map[string]info.UserAndBlog {
	return make(map[string]info.UserAndBlog)
}

// ListUserInfosByUserIds returns user info for a list of user IDs
func (this *UserService) ListUserInfosByUserIds(userIds []string) []info.User {
	if len(userIds) == 0 {
		return []info.User{}
	}

	users := []info.User{}
	query := "SELECT id, email, username, logo FROM users WHERE id = ANY($1)"
	rows, err := db.DB.Query(query, pq.Array(userIds))
	if err != nil {
		Log(err.Error())
		return []info.User{}
	}
	defer rows.Close()

	for rows.Next() {
		var user info.User
		err := rows.Scan(&user.UserId, &user.Email, &user.Username, &user.Logo)
		if err != nil {
			Log(err.Error())
			continue
		}
		users = append(users, user)
	}
	return users
}
