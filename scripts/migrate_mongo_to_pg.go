package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	mongoURL      = "mongodb://localhost:27017/leanote"
	postgresURL   = "host=localhost port=5432 user=leanote password=your_password dbname=leanote sslmode=disable"
	mongoDatabase = "leanote"
)

var (
	mongoSession *mgo.Session
	postgresDB   *sql.DB
)

func main() {
	// Connect to MongoDB
	var err error
	mongoSession, err = mgo.Dial(mongoURL)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoSession.Close()

	mongoDB := mongoSession.DB(mongoDatabase)

	// Connect to PostgreSQL
	postgresDB, err = sql.Open("postgres", postgresURL)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer postgresDB.Close()

	err = postgresDB.Ping()
	if err != nil {
		log.Fatalf("Failed to ping PostgreSQL: %v", err)
	}

	fmt.Println("Starting migration from MongoDB to PostgreSQL...")
	fmt.Println("Connected to both databases successfully")

	// Migrate in order of dependencies
	fmt.Println("\n=== Migrating Users ===")
	migrateUsers(mongoDB)

	fmt.Println("\n=== Migrating Notebooks ===")
	migrateNotebooks(mongoDB)

	fmt.Println("\n=== Migrating Notes ===")
	migrateNotes(mongoDB)

	fmt.Println("\n=== Migrating Note Contents ===")
	migrateNoteContents(mongoDB)

	fmt.Println("\n=== Migrating Tags ===")
	migrateTags(mongoDB)

	fmt.Println("\n=== Migrating Attachments ===")
	migrateAttachments(mongoDB)

	fmt.Println("\n=== Migration Complete ===")
}

type MongoUser struct {
	Id              bson.ObjectId `bson:"_id"`
	Email           string        `bson:"Email"`
	Verified        bool          `bson:"Verified"`
	Username        string        `bson:"Username"`
	UsernameRaw     string        `bson:"UsernameRaw"`
	Pwd             string        `bson:"Pwd"`
	CreatedTime     string        `bson:"CreatedTime"`
	Logo            string        `bson:"Logo"`
	Theme           string        `bson:"Theme"`
	NotebookWidth   int           `bson:"NotebookWidth"`
	NoteListWidth   int           `bson:"NoteListWidth"`
	MdEditorWidth   int           `bson:"MdEditorWidth"`
	LeftIsMin       bool          `bson:"LeftIsMin"`
	ThirdUserId     string        `bson:"ThirdUserId"`
	ThirdUsername   string        `bson:"ThirdUsername"`
	ThirdType       int           `bson:"ThirdType"`
	ImageNum        int           `bson:"ImageNum"`
	ImageSize       int           `bson:"ImageSize"`
	AttachNum       int           `bson:"AttachNum"`
	AttachSize      int           `bson:"AttachSize"`
	FromUserId      bson.ObjectId  `bson:"FromUserId"`
	AccountType     string        `bson:"AccountType"`
	AccountStartTime string       `bson:"AccountStartTime"`
	AccountEndTime   string       `bson:"AccountEndTime"`
	MaxImageNum      int          `bson:"MaxImageNums"`
	MaxImageSize     int          `bson:"MaxImageSize"`
	MaxAttachNum     int          `bson:"MaxAttachNum"`
	MaxAttachSize    int          `bson:"MaxAttachSize"`
	MaxPerAttachSize int          `bson:"MaxPerAttachSize"`
	Usn              int          `bson:"Usn"`
	FullSyncBefore   string       `bson:"FullSyncBefore"`
	IsDeleted        bool         `bson:"IsDeleted"`
}

func migrateUsers(db *mgo.Database) {
	collection := db.C("users")

	var users []MongoUser
	err := collection.Find(nil).All(&users)
	if err != nil {
		log.Printf("Error fetching users: %v", err)
		return
	}

	successCount := 0
	for _, user := range users {
		userId := objectIdToUUID(user.Id.Hex())
		fromUserId := ""
		if user.FromUserId != "" {
			fromUserId = objectIdToUUID(user.FromUserId.Hex())
		}

		query := `INSERT INTO users (
			id, email, verified, username, username_raw, pwd, created_time,
			logo, theme, notebook_width, note_list_width, md_editor_width, left_is_min,
			third_user_id, third_username, third_type, image_num, image_size,
			attach_num, attach_size, from_user_id, account_type, account_start_time,
			account_end_time, max_image_num, max_image_size, max_attach_num,
			max_attach_size, max_per_attach_size, usn, full_sync_before, is_deleted
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29,
			$30, $31)`

		_, err := postgresDB.Exec(query,
			userId,
			user.Email,
			user.Verified,
			user.Username,
			user.UsernameRaw,
			user.Pwd,
			parseTime(user.CreatedTime),
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
			fromUserId,
			user.AccountType,
			parseTime(user.AccountStartTime),
			parseTime(user.AccountEndTime),
			user.MaxImageNum,
			user.MaxImageSize,
			user.MaxAttachNum,
			user.MaxAttachSize,
			user.MaxPerAttachSize,
			user.Usn,
			parseTime(user.FullSyncBefore),
			user.IsDeleted,
		)

		if err != nil {
			log.Printf("Error inserting user %s: %v", user.Id, err)
		} else {
			successCount++
		}
	}

	fmt.Printf("Migrated %d/%d users\n", successCount, len(users))
}

type MongoNotebook struct {
	Id               bson.ObjectId `bson:"_id"`
	UserId           bson.ObjectId `bson:"UserId"`
	ParentNotebookId bson.ObjectId `bson:"ParentNotebookId"`
	Seq              int           `bson:"Seq"`
	Title            string        `bson:"Title"`
	UrlTitle         string        `bson:"UrlTitle"`
	NumberNotes      int           `bson:"NumberNotes"`
	IsTrash          bool          `bson:"IsTrash"`
	IsBlog           bool          `bson:"IsBlog"`
	CreatedTime      string        `bson:"CreatedTime"`
	UpdatedTime      string        `bson:"UpdatedTime"`
	Usn              int           `bson:"Usn"`
	IsDeleted        bool          `bson:"IsDeleted"`
}

func migrateNotebooks(db *mgo.Database) {
	collection := db.C("notebooks")

	var notebooks []MongoNotebook
	err := collection.Find(nil).All(&notebooks)
	if err != nil {
		log.Printf("Error fetching notebooks: %v", err)
		return
	}

	successCount := 0
	for _, notebook := range notebooks {
		notebookId := objectIdToUUID(notebook.Id.Hex())
		userId := objectIdToUUID(notebook.UserId.Hex())
		parentNotebookId := ""
		if notebook.ParentNotebookId != "" {
			parentNotebookId = objectIdToUUID(notebook.ParentNotebookId.Hex())
		}

		query := `INSERT INTO notebooks (
			id, user_id, parent_notebook_id, seq, title, url_title, number_notes,
			is_trash, is_blog, created_time, updated_time, usn, is_deleted
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

		_, err := postgresDB.Exec(query,
			notebookId,
			userId,
			parentNotebookId,
			notebook.Seq,
			notebook.Title,
			notebook.UrlTitle,
			notebook.NumberNotes,
			notebook.IsTrash,
			notebook.IsBlog,
			parseTime(notebook.CreatedTime),
			parseTime(notebook.UpdatedTime),
			notebook.Usn,
			notebook.IsDeleted,
		)

		if err != nil {
			log.Printf("Error inserting notebook %s: %v", notebook.Id, err)
		} else {
			successCount++
		}
	}

	fmt.Printf("Migrated %d/%d notebooks\n", successCount, len(notebooks))
}

type MongoNote struct {
	Id            bson.ObjectId   `bson:"_id"`
	UserId        bson.ObjectId   `bson:"UserId"`
	CreatedUserId bson.ObjectId   `bson:"CreatedUserId"`
	NotebookId    bson.ObjectId   `bson:"NotebookId"`
	Title         string          `bson:"Title"`
	Desc          string          `bson:"Desc"`
	Src           string          `bson:"Src"`
	ImgSrc        string          `bson:"ImgSrc"`
	Tags          []string        `bson:"Tags"`
	IsTrash       bool            `bson:"IsTrash"`
	IsBlog        bool            `bson:"IsBlog"`
	UrlTitle      string          `bson:"UrlTitle"`
	IsRecommend   bool            `bson:"IsRecommend"`
	IsTop         bool            `bson:"IsTop"`
	HasSelfDefined bool           `bson:"HasSelfDefined"`
	ReadNum       int             `bson:"ReadNum"`
	LikeNum       int             `bson:"LikeNum"`
	CommentNum    int             `bson:"CommentNum"`
	IsMarkdown    bool            `bson:"IsMarkdown"`
	AttachNum     int             `bson:"AttachNum"`
	CreatedTime   string          `bson:"CreatedTime"`
	UpdatedTime   string          `bson:"UpdatedTime"`
	RecommendTime string          `bson:"RecommendTime"`
	PublicTime    string          `bson:"PublicTime"`
	UpdatedUserId bson.ObjectId   `bson:"UpdatedUserId"`
	Usn           int             `bson:"Usn"`
	IsDeleted     bool            `bson:"IsDeleted"`
}

func migrateNotes(db *mgo.Database) {
	collection := db.C("notes")

	var notes []MongoNote
	err := collection.Find(nil).All(&notes)
	if err != nil {
		log.Printf("Error fetching notes: %v", err)
		return
	}

	successCount := 0
	for _, note := range notes {
		noteId := objectIdToUUID(note.Id.Hex())
		userId := objectIdToUUID(note.UserId.Hex())
		createdUserId := ""
		if note.CreatedUserId != "" {
			createdUserId = objectIdToUUID(note.CreatedUserId.Hex())
		}
		notebookId := objectIdToUUID(note.NotebookId.Hex())
		updatedUserId := ""
		if note.UpdatedUserId != "" {
			updatedUserId = objectIdToUUID(note.UpdatedUserId.Hex())
		}

		tagsJSON, _ := json.Marshal(note.Tags)

		query := `INSERT INTO notes (
			id, user_id, created_user_id, notebook_id, title, description, src, img_src,
			tags, is_trash, is_blog, url_title, is_recommend, is_top, has_self_defined,
			read_num, like_num, comment_num, is_markdown, attach_num, created_time,
			updated_time, recommend_time, public_time, updated_user_id, usn, is_deleted
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27)`

		_, err := postgresDB.Exec(query,
			noteId,
			userId,
			createdUserId,
			notebookId,
			note.Title,
			note.Desc,
			note.Src,
			note.ImgSrc,
			tagsJSON,
			note.IsTrash,
			note.IsBlog,
			note.UrlTitle,
			note.IsRecommend,
			note.IsTop,
			note.HasSelfDefined,
			note.ReadNum,
			note.LikeNum,
			note.CommentNum,
			note.IsMarkdown,
			note.AttachNum,
			parseTime(note.CreatedTime),
			parseTime(note.UpdatedTime),
			parseTime(note.RecommendTime),
			parseTime(note.PublicTime),
			updatedUserId,
			note.Usn,
			note.IsDeleted,
		)

		if err != nil {
			log.Printf("Error inserting note %s: %v", note.Id, err)
		} else {
			successCount++
		}
	}

	fmt.Printf("Migrated %d/%d notes\n", successCount, len(notes))
}

type MongoNoteContent struct {
	NoteId       bson.ObjectId `bson:"_id"`
	UserId       bson.ObjectId `bson:"UserId"`
	IsBlog       bool          `bson:"IsBlog"`
	Content      string        `bson:"Content"`
	Abstract     string        `bson:"Abstract"`
	CreatedTime  string        `bson:"CreatedTime"`
	UpdatedTime  string        `bson:"UpdatedTime"`
	UpdatedUserId bson.ObjectId `bson:"UpdatedUserId"`
}

func migrateNoteContents(db *mgo.Database) {
	collection := db.C("note_contents")

	var noteContents []MongoNoteContent
	err := collection.Find(nil).All(&noteContents)
	if err != nil {
		log.Printf("Error fetching note contents: %v", err)
		return
	}

	successCount := 0
	for _, nc := range noteContents {
		noteId := objectIdToUUID(nc.NoteId.Hex())
		userId := objectIdToUUID(nc.UserId.Hex())
		updatedUserId := ""
		if nc.UpdatedUserId != "" {
			updatedUserId = objectIdToUUID(nc.UpdatedUserId.Hex())
		}

		query := `INSERT INTO note_contents (
			note_id, user_id, is_blog, content, abstract, created_time,
			updated_time, updated_user_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

		_, err := postgresDB.Exec(query,
			noteId,
			userId,
			nc.IsBlog,
			nc.Content,
			nc.Abstract,
			parseTime(nc.CreatedTime),
			parseTime(nc.UpdatedTime),
			updatedUserId,
		)

		if err != nil {
			log.Printf("Error inserting note content %s: %v", nc.NoteId, err)
		} else {
			successCount++
		}
	}

	fmt.Printf("Migrated %d/%d note contents\n", successCount, len(noteContents))
}

type MongoTag struct {
	Id         bson.ObjectId `bson:"_id"`
	UserId     bson.ObjectId `bson:"UserId"`
	Tag        string        `bson:"Tag"`
	Usn        int           `bson:"Usn"`
	Count      int           `bson:"Count"`
	CreatedTime string        `bson:"CreatedTime"`
	UpdatedTime string        `bson:"UpdatedTime"`
	IsDeleted  bool          `bson:"IsDeleted"`
}

func migrateTags(db *mgo.Database) {
	collection := db.C("tags")

	var tags []MongoTag
	err := collection.Find(nil).All(&tags)
	if err != nil {
		log.Printf("Error fetching tags: %v", err)
		return
	}

	successCount := 0
	for _, tag := range tags {
		tagId := objectIdToUUID(tag.Id.Hex())
		userId := objectIdToUUID(tag.UserId.Hex())

		query := `INSERT INTO tags (
			id, user_id, tag, usn, count, created_time, updated_time, is_deleted
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

		_, err := postgresDB.Exec(query,
			tagId,
			userId,
			tag.Tag,
			tag.Usn,
			tag.Count,
			parseTime(tag.CreatedTime),
			parseTime(tag.UpdatedTime),
			tag.IsDeleted,
		)

		if err != nil {
			log.Printf("Error inserting tag %s: %v", tag.Id, err)
		} else {
			successCount++
		}
	}

	fmt.Printf("Migrated %d/%d tags\n", successCount, len(tags))
}

type MongoAttach struct {
	Id           bson.ObjectId `bson:"_id"`
	NoteId       bson.ObjectId `bson:"NoteId"`
	UploadUserId bson.ObjectId `bson:"UploadUserId"`
	Name         string        `bson:"Name"`
	Title        string        `bson:"Title"`
	Size         int64         `bson:"Size"`
	Type         string        `bson:"Type"`
	Path         string        `bson:"Path"`
	CreatedTime  string        `bson:"CreatedTime"`
}

func migrateAttachments(db *mgo.Database) {
	collection := db.C("attachs")

	var attaches []MongoAttach
	err := collection.Find(nil).All(&attaches)
	if err != nil {
		log.Printf("Error fetching attachments: %v", err)
		return
	}

	successCount := 0
	for _, attach := range attaches {
		attachId := objectIdToUUID(attach.Id.Hex())
		noteId := objectIdToUUID(attach.NoteId.Hex())
		uploadUserId := objectIdToUUID(attach.UploadUserId.Hex())

		query := `INSERT INTO attachs (
			id, note_id, upload_user_id, name, title, size, type, path, created_time
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

		_, err := postgresDB.Exec(query,
			attachId,
			noteId,
			uploadUserId,
			attach.Name,
			attach.Title,
			attach.Size,
			attach.Type,
			attach.Path,
			parseTime(attach.CreatedTime),
		)

		if err != nil {
			log.Printf("Error inserting attachment %s: %v", attach.Id, err)
		} else {
			successCount++
		}
	}

	fmt.Printf("Migrated %d/%d attachments\n", successCount, len(attaches))
}

// Helper function to convert MongoDB ObjectId to UUID
// This is a simple conversion - in production you might want to preserve the original ObjectId
// or use a more sophisticated mapping strategy
func objectIdToUUID(objectId string) string {
	// For now, generate a new UUID
	// In production, you might want to store a mapping or use a UUID v5 based on the ObjectId
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		objectId[0:8],
		objectId[8:12],
		objectId[12:16],
		objectId[16:20],
		objectId[20:24])
}

// Helper function to parse MongoDB time strings
func parseTime(timeStr string) interface{} {
	if timeStr == "" || timeStr == "0001-01-01T00:00:00Z" {
		return nil
	}
	return timeStr
}
