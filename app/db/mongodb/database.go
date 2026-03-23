package mongodb

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	. "github.com/leanote/leanote/app/lea"

	"github.com/leanote/leanote/app/db"
)

func (m *MongoDatabase) Insert(table string, data interface{}) bool {
	collection := m.getCollection(table)
	if collection == nil {
		Log("Unknown table: " + table)
		return false
	}

	err := collection.Insert(data)
	if err != nil {
		Log("Insert error: " + err.Error())
		return false
	}
	return true
}

func (m *MongoDatabase) Update(table string, id string, data interface{}) bool {
	collection := m.getCollection(table)
	if collection == nil {
		return false
	}

	err := collection.UpdateId(bson.ObjectIdHex(id), data)
	if err != nil {
		Log("Update error: " + err.Error())
		return false
	}
	return true
}

func (m *MongoDatabase) Delete(table string, id string) bool {
	collection := m.getCollection(table)
	if collection == nil {
		return false
	}

	err := collection.RemoveId(bson.ObjectIdHex(id))
	if err != nil {
		Log("Delete error: " + err.Error())
		return false
	}
	return true
}

func (m *MongoDatabase) getCollection(table string) *mgo.Collection {
	switch table {
	case "users":
		return m.Users
	case "notebooks":
		return m.Notebooks
	case "notes":
		return m.Notes
	case "note_contents":
		return m.NoteContents
	case "note_content_histories":
		return m.NoteContentHistories
	case "share_notes":
		return m.ShareNotes
	case "share_notebooks":
		return m.ShareNotebooks
	case "has_share_notes":
		return m.HasShareNotes
	case "blogs":
		return m.Blogs
	case "groups":
		return m.Groups
	case "group_users":
		return m.GroupUsers
	case "tags":
		return m.Tags
	case "note_tags":
		return m.NoteTags
	case "tag_count":
		return m.TagCounts
	case "user_blogs":
		return m.UserBlogs
	case "tokens":
		return m.Tokens
	case "suggestions":
		return m.Suggestions
	case "albums":
		return m.Albums
	case "files":
		return m.Files
	case "attachs":
		return m.Attachs
	case "note_images":
		return m.NoteImages
	case "configs":
		return m.Configs
	case "email_logs":
		return m.EmailLogs
	case "blog_likes":
		return m.BlogLikes
	case "blog_comments":
		return m.BlogComments
	case "reports":
		return m.Reports
	case "blog_singles":
		return m.BlogSingles
	case "themes":
		return m.Themes
	case "sessions":
		return m.Sessions
	default:
		return nil
	}
}

func (m *MongoDatabase) InsertOld(collection interface{}, data interface{}) bool {
	if coll, ok := collection.(*mgo.Collection); ok {
		err := coll.Insert(data)
		return err == nil
	}
	return false
}

func (m *MongoDatabase) UpdateByIdAndUserId(collection interface{}, id, userId string, data interface{}) bool {
	if coll, ok := collection.(*mgo.Collection); ok {
		err := coll.Update(bson.M{"_id": bson.ObjectIdHex(id), "UserId": bson.ObjectIdHex(userId)}, data)
		return err == nil
	}
	return false
}

func (m *MongoDatabase) GetByIdAndUserId(collection interface{}, id, userId string, data interface{}) {
	if coll, ok := collection.(*mgo.Collection); ok {
		coll.Find(bson.M{"_id": bson.ObjectIdHex(id), "UserId": bson.ObjectIdHex(userId)}).One(data)
	}
}

func (m *MongoDatabase) GetByQ(collection interface{}, query interface{}, data interface{}) {
	if coll, ok := collection.(*mgo.Collection); ok {
		coll.Find(query).One(data)
	}
}

func (m *MongoDatabase) ListByQ(collection interface{}, query interface{}, data interface{}) {
	if coll, ok := collection.(*mgo.Collection); ok {
		coll.Find(query).All(data)
	}
}

func (m *MongoDatabase) Count(collection interface{}, query interface{}) int {
	if coll, ok := collection.(*mgo.Collection); ok {
		count, _ := coll.Find(query).Count()
		return count
	}
	return 0
}

func (m *MongoDatabase) Has(collection interface{}, query interface{}) bool {
	return m.Count(collection, query) > 0
}

func (m *MongoDatabase) BatchInsert(table string, data []interface{}) (int, error) {
	collection := m.getCollection(table)
	if collection == nil || len(data) == 0 {
		return 0, nil
	}

	bulk := collection.Bulk()
	for _, item := range data {
		bulk.Insert(item)
	}

	_, err := bulk.Run()
	if err != nil {
		return 0, err
	}
	return len(data), nil
}

func (m *MongoDatabase) BatchUpdate(table string, ids []string, data interface{}) (int, error) {
	collection := m.getCollection(table)
	if collection == nil || len(ids) == 0 {
		return 0, nil
	}

	successCount := 0
	for _, id := range ids {
		if m.Update(table, id, data) {
			successCount++
		}
	}
	return successCount, nil
}

func (m *MongoDatabase) BatchDelete(table string, ids []string) (int, error) {
	collection := m.getCollection(table)
	if collection == nil || len(ids) == 0 {
		return 0, nil
	}

	objectIds := make([]bson.ObjectId, len(ids))
	for i, id := range ids {
		objectIds[i] = bson.ObjectIdHex(id)
	}

	_, err := collection.RemoveAll(bson.M{"_id": bson.M{"$in": objectIds}})
	if err != nil {
		return 0, err
	}
	return len(ids), nil
}

func (m *MongoDatabase) Select(table string, fields []string, where string, args ...interface{}) ([]map[string]interface{}, error) {
	collection := m.getCollection(table)
	if collection == nil {
		return nil, nil
	}

	selector := bson.M{}
	if fields != nil && len(fields) > 0 {
		selector = bson.M{}
		for _, field := range fields {
			selector[field] = 1
		}
	}

	var results []map[string]interface{}
	err := collection.Find(selector).All(&results)
	return results, err
}

func (m *MongoDatabase) Join(mainTable, joinTable, joinType, onCondition string, where string, args ...interface{}) ([]map[string]interface{}, error) {
	return nil, nil
}

func (m *MongoDatabase) Paginate(table string, page, pageSize int, where string, orderBy string, args ...interface{}) (db.PaginationResult, error) {
	result := db.PaginationResult{
		Page:     page,
		PageSize: pageSize,
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	collection := m.getCollection(table)
	if collection == nil {
		return result, nil
	}

	total, _ := collection.Find(nil).Count()
	result.Total = int64(total)
	result.TotalPages = (total + pageSize - 1) / pageSize

	query := collection.Find(nil).Skip((page - 1) * pageSize).Limit(pageSize)

	var data []map[string]interface{}
	query.All(&data)
	result.Data = data

	return result, nil
}

func (m *MongoDatabase) Begin() (db.Transaction, error) {
	session := m.session.Copy()
	return NewMongoTransaction(session), nil
}
