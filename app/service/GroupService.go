package service

import (
	"github.com/leanote/leanote/app/db"
	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
	"github.com/lib/pq"
	"time"
)

// 用户组, 用户组用户管理

type GroupService struct {
}

// 添加分组
func (this *GroupService) AddGroup(userId, title string) (bool, info.Group) {
	group := info.Group{
		GroupId:     db.NewUUID(),
		UserId:      userId,
		Title:       title,
		CreatedTime: time.Now(),
	}
	return db.Insert(db.Groups, group), group
}

// 删除分组
// 判断是否有好友
func (this *GroupService) DeleteGroup(userId, groupId string) (ok bool, msg string) {
	/*
		// 简化版本
	*/
	if !this.isMyGroup(userId, groupId) {
		return false, "notMyGroup"
	}

	// 删除分组后, 需要删除所有用户分享到该组的笔记本, 笔记

	shareService.DeleteAllShareNotebookGroup(groupId)
	shareService.DeleteAllShareNoteGroup(groupId)

	query := `DELETE FROM group_users WHERE group_id = $1`
	db.DB.Exec(query, groupId)
	return db.DeleteByIdAndUserId(db.Groups, groupId, userId), ""

	// TODO 删除分组后, 在shareNote, shareNotebook中也要删除
}

// 修改group标题
func (this *GroupService) UpdateGroupTitle(userId, groupId, title string) (ok bool) {
	return db.UpdateByIdAndUserIdField(db.Groups, groupId, userId, "Title", title)
}

// 得到用户的所有分组(包括下的所有用户)
func (this *GroupService) GetGroupsAndUsers(userId string) []info.Group {
	/*
			// 得到我的分组
			groups := []info.Group{}
		query := `SELECT group_id, user_id, title, created_time FROM groups WHERE user_id = $1`
		rows, err := db.DB.Query(query, userId)
		if err != nil {
			return groups
		}
		defer rows.Close()

		for rows.Next() {
			var group info.Group
			err := rows.Scan(&group.GroupId, &group.UserId, &group.Title, &group.CreatedTime)
			if err == nil {
				groups = append(groups, group)
			}
		}
	*/
	// 我的分组, 及我所属的分组
	groups := this.GetGroupsContainOf(userId)

	// 得到其下的用户
	for i, group := range groups {
		group.Users = this.GetUsers(group.GroupId)
		groups[i] = group
	}
	return groups
}

// 仅仅得到所有分组
func (this *GroupService) GetGroups(userId string) []info.Group {
	// 得到分组s
	groups := []info.Group{}
	query := `SELECT group_id, user_id, title, created_time FROM groups WHERE user_id = $1`
	rows, err := db.DB.Query(query, userId)
	if err != nil {
		return groups
	}
	defer rows.Close()

	for rows.Next() {
		var group info.Group
		err := rows.Scan(&group.GroupId, &group.UserId, &group.Title, &group.CreatedTime)
		if err == nil {
			groups = append(groups, group)
		}
	}
	return groups
}

// 得到我的和我所属组的ids
func (this *GroupService) GetMineAndBelongToGroupIds(userId string) []string {
	// 所属组
	groupIds := this.GetBelongToGroupIds(userId)

	m := map[string]bool{}
	for _, groupId := range groupIds {
		m[groupId] = true
	}

	// 我的组
	myGroups := this.GetGroups(userId)

	for _, group := range myGroups {
		if !m[group.GroupId] {
			groupIds = append(groupIds, group.GroupId)
		}
	}

	return groupIds
}

// 获取包含此用户的组对象数组
// 获取该用户所属组, 和我的组
func (this *GroupService) GetGroupsContainOf(userId string) []info.Group {
	// 我的组
	myGroups := this.GetGroups(userId)
	myGroupMap := map[string]bool{}

	for _, group := range myGroups {
		myGroupMap[group.GroupId] = true
	}

	// 所属组
	groupIds := this.GetBelongToGroupIds(userId)
	if len(groupIds) == 0 {
		return myGroups
	}

	// Build PostgreSQL query with IN clause
	query := "SELECT * FROM groups WHERE id = ANY($1)"
	rows, err := db.DB.Query(query, pq.Array(groupIds))
	if err != nil {
		Log(err.Error())
		return myGroups
	}
	defer rows.Close()

	for rows.Next() {
		var group info.Group
		err := rows.Scan(&group.GroupId, &group.UserId, &group.Title, &group.CreatedTime)
		if err != nil {
			Log(err.Error())
			continue
		}
		if !myGroupMap[group.GroupId] {
			myGroups = append(myGroups, group)
		}
	}

	return myGroups
}

// 得到分组, shareService用
func (this *GroupService) GetGroup(userId, groupId string) info.Group {
	// 得到分组s
	group := info.Group{}
	db.GetByIdAndUserId(db.Groups, groupId, userId, &group)
	return group
}

// 得到某分组下的用户
func (this *GroupService) GetUsers(groupId string) []info.User {
	// 得到UserIds
	query := "SELECT user_id FROM group_users WHERE group_id = $1"
	rows, err := db.DB.Query(query, groupId)
	if err != nil {
		Log(err.Error())
		return nil
	}
	defer rows.Close()

	userIds := []string{}
	for rows.Next() {
		var userId string
		err := rows.Scan(&userId)
		if err != nil {
			Log(err.Error())
			continue
		}
		userIds = append(userIds, userId)
	}
	if len(userIds) == 0 {
		return nil
	}
	// 得到userInfos
	return userService.ListUserInfosByUserIds(userIds)
}

// 得到我所属的所有分组ids
func (this *GroupService) GetBelongToGroupIds(userId string) []string {
	// 得到UserIds
	query := "SELECT group_id FROM group_users WHERE user_id = $1"
	rows, err := db.DB.Query(query, userId)
	if err != nil {
		Log(err.Error())
		return nil
	}
	defer rows.Close()

	groupIds := []string{}
	for rows.Next() {
		var groupUser info.GroupUser
		err := rows.Scan(&groupUser.GroupId)
		if err != nil {
			Log(err.Error())
			continue
		}
		groupIds = append(groupIds, groupUser.GroupId)
	}
	return groupIds
}

func (this *GroupService) isMyGroup(ownUserId, groupId string) (ok bool) {
	query := "SELECT COUNT(*) FROM groups WHERE id = $1 AND user_id = $2"
	var count int
	err := db.DB.QueryRow(query, groupId, ownUserId).Scan(&count)
	if err != nil {
		Log(err.Error())
		return false
	}
	return count > 0
}

// 判断组中是否包含指定用户
func (this *GroupService) IsExistsGroupUser(userId, groupId string) (ok bool) {
	// 如果我拥有这个组, 那也行
	if this.isMyGroup(userId, groupId) {
		return true
	}
	query := "SELECT COUNT(*) FROM group_users WHERE user_id = $1 AND group_id = $2"
	var count int
	err := db.DB.QueryRow(query, userId, groupId).Scan(&count)
	if err != nil {
		Log(err.Error())
		return false
	}
	return count > 0
}

// 为group添加用户
// 用户是否已存在?
func (this *GroupService) AddUser(ownUserId, groupId, userId string) (ok bool, msg string) {
	// groupId是否是ownUserId的?
	/*
		if !this.IsExistsGroupUser(ownUserId, groupId) {
			return false, "forbiddenNotMyGroup"
		}
	*/
	if !this.isMyGroup(ownUserId, groupId) {
		return false, "forbiddenNotMyGroup"
	}

	// 是否已存在
	query := "SELECT COUNT(*) FROM group_users WHERE group_id = $1 AND user_id = $2"
	var count int
	err := db.DB.QueryRow(query, groupId, userId).Scan(&count)
	if err != nil {
		Log(err.Error())
		return false, "errorCheckingUser"
	}
	if count > 0 {
		return false, "userExistsInGroup"
	}

	groupUser := info.GroupUser{
		GroupUserId: db.NewUUID(),
		GroupId:     groupId,
		UserId:      userId,
		CreatedTime: time.Now(),
	}
	query = "INSERT INTO group_users (id, group_id, user_id, created_time) VALUES ($1, $2, $3, $4)"
	_, err = db.DB.Exec(query, groupUser.GroupUserId, groupUser.GroupId, groupUser.UserId, groupUser.CreatedTime)
	return err == nil, ""
}

// 删除用户
func (this *GroupService) DeleteUser(ownUserId, groupId, userId string) (ok bool, msg string) {
	// groupId是否是ownUserId的?
	/*
		if !this.IsExistsGroupUser(ownUserId, groupId) {
			return false, "forbiddenNotMyGroup"
		}
	*/
	if !this.isMyGroup(ownUserId, groupId) {
		return false, "forbiddenNotMyGroup"
	}

	// 删除该用户分享到本组的笔记本, 笔记
	shareService.DeleteShareNotebookGroupWhenDeleteGroupUser(userId, groupId)
	shareService.DeleteShareNoteGroupWhenDeleteGroupUser(userId, groupId)

	query := "DELETE FROM group_users WHERE group_id = $1 AND user_id = $2"
	_, err := db.DB.Exec(query, groupId, userId)
	return err == nil, ""
}
