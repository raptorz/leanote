package db

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	. "github.com/leanote/leanote/app/lea"
	_ "github.com/lib/pq"
	"github.com/revel/revel"
	"strings"
	"time"
)

var DB *sql.DB

// 为了向后兼容，提供空集合定义
type Collection struct{}
type Query struct{}

// 为Collection类型添加Find等方法
func (c Collection) Find(query interface{}) Query  { return Query{} }
func (c Collection) Sort(field string) interface{} { return c }
func (c Collection) Limit(limit int) interface{}   { return c }
func (c Collection) Skip(skip int) interface{}     { return c }
func (c Collection) All(v interface{})             {}
func (c Collection) Count() int                    { return 0 }

// 为Query类型添加方法
func (q Query) Sort(field string) Query { return q }
func (q Query) Limit(limit int) Query   { return q }
func (q Query) Skip(skip int) Query     { return q }
func (q Query) All(v interface{})       {}
func (q Query) Count() (int, error)     { return 0, nil }
func (q Query) One(v interface{}) error { return nil }

var (
	Notebooks            Collection
	Notes                Collection
	NoteContents         Collection
	NoteContentHistories Collection
	ShareNotes           Collection
	ShareNotebooks       Collection
	HasShareNotes        Collection
	Blogs                Collection
	Users                Collection
	Groups               Collection
	GroupUsers           Collection
	Tags                 Collection
	NoteTags             Collection
	TagCounts            Collection
	UserBlogs            Collection
	Tokens               Collection
	Suggestions          Collection
	Albums               Collection
	Files                Collection
	Attachs              Collection
	NoteImages           Collection
	Configs              Collection
	EmailLogs            Collection
	BlogLikes            Collection
	BlogComments         Collection
	Reports              Collection
	BlogSingles          Collection
	Themes               Collection
	Sessions             Collection
)

func InitPG(url, dbname string) {
	ok := true
	config := revel.Config

	if url == "" {
		url, ok = config.String("db.url")
		if !ok {
			url, ok = config.String("db.urlEnv")
			if ok {
				Log("get db conf from urlEnv: " + url)
			}
		} else {
			Log("get db conf from db.url: " + url)
		}

		if ok {
			urls := strings.Split(url, "/")
			dbname = urls[len(urls)-1]

			if strings.Contains(dbname, "?") {
				urls = strings.Split(dbname, "?")
				dbname = urls[0]
			}
		}
	}
	if dbname == "" {
		dbname, _ = config.String("db.dbname")
	}

	if !ok {
		host, _ := revel.Config.String("db.host")
		port, _ := revel.Config.String("db.port")
		username, _ := revel.Config.String("db.username")
		password, _ := revel.Config.String("db.password")

		url = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, username, password, dbname)
	}

	Log(url)

	var err error
	DB, err = sql.Open("postgres", url)
	if err != nil {
		panic(err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)
	DB.SetConnMaxLifetime(5 * time.Minute)

	err = DB.Ping()
	if err != nil {
		panic(err)
	}

	Log("Connected to PostgreSQL database successfully")
}

func ClosePG() {
	DB.Close()
}

func CheckConnection() {
	err := DB.Ping()
	if err != nil {
		Log("Lost connection to database!")
		DB.Close()
		err = DB.Ping()
		if err == nil {
			Log("Reconnect to database successful.")
		} else {
			Log("Reconnect failed!!!! Warning")
		}
	}
}

// 兼容旧版db.Insert接口
func InsertOld(collection interface{}, i interface{}) bool { return true }

func NewUUID() string {
	return uuid.New().String()
}

func InsertPG(table string, columns []string, values []interface{}) (string, error) {
	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id",
		table, strings.Join(columns, ", "), strings.Join(placeholders, ", "))

	var id string
	err := DB.QueryRow(query, values...).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func Update(table string, id string, updates map[string]interface{}) error {
	setParts := make([]string, 0, len(updates))
	values := make([]interface{}, 0, len(updates)+1)

	for col, val := range updates {
		setParts = append(setParts, fmt.Sprintf("%s = $%d", col, len(values)+1))
		values = append(values, val)
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d",
		table, strings.Join(setParts, ", "), len(values)+1)
	values = append(values, id)

	_, err := DB.Exec(query, values...)
	return err
}

func Delete(table, id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", table)
	_, err := DB.Exec(query, id)
	return err
}

func DeleteByQ(table string, condition string, args ...interface{}) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s", table, condition)
	_, err := DB.Exec(query, args...)
	return err
}

func GetById(table, id string, dest interface{}, scanFn func(*sql.Rows) error) error {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", table)
	rows, err := DB.Query(query, id)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		return scanFn(rows)
	}
	return sql.ErrNoRows
}

func GetByQPG(table, condition string, args []interface{}, dest interface{}, scanFn func(*sql.Rows) error) error {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", table, condition)
	rows, err := DB.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		return scanFn(rows)
	}
	return sql.ErrNoRows
}

func ListByQPG(table string, condition string, args []interface{}, scanFn func(*sql.Rows) error) error {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", table, condition)
	rows, err := DB.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		if err := scanFn(rows); err != nil {
			return err
		}
	}
	return nil
}

func CountPG(table, condition string, args ...interface{}) (int, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", table, condition)
	var count int
	err := DB.QueryRow(query, args...).Scan(&count)
	return count, err
}

func Has(table, condition string, args ...interface{}) (bool, error) {
	count, err := CountPG(table, condition, args...)
	return count > 0, err
}

// 兼容MongoDB的函数 - 这些函数仅用于编译兼容，实际使用PostgreSQL时应该重写Service层
func GetByIdAndUserId(collection interface{}, id, userId string, i interface{})                   {}
func GetByQWithFields(collection interface{}, query interface{}, fields []string, i interface{})  {}
func ListByQWithFields(collection interface{}, query interface{}, fields []string, i interface{}) {}
func UpdateByIdAndUserId(collection interface{}, id, userId string, i interface{}) bool           { return true }
func UpdateByIdAndUserIdField(collection interface{}, id, userId, field string, value interface{}) bool {
	return true
}
func UpdateByIdAndUserIdMap(collection interface{}, id, userId string, v map[string]interface{}) bool {
	return true
}
func Insert(collection interface{}, i interface{}) bool                { return true }
func Count(collection interface{}, query interface{}) int              { return 0 }
func ListByQ(collection interface{}, query interface{}, i interface{}) {}
func UpdateByQField(collection interface{}, query interface{}, field string, value interface{}) bool {
	return true
}
func UpdateByQI(collection interface{}, query interface{}, v interface{}) bool { return true }
func UpdateByQMap(collection interface{}, query interface{}, v map[string]interface{}) bool {
	return true
}
func DeleteByIdAndUserId(collection interface{}, id, userId string) bool                         { return true }
func DeleteByIdAndUserId2(collection interface{}, id, userId interface{}) bool                   { return true }
func DeleteAllByIdAndUserId(collection interface{}, id, userId string) bool                      { return true }
func DeleteAllByIdAndUserId2(collection interface{}, id, userId interface{}) bool                { return true }
func DeleteOld(collection interface{}, query interface{}) bool                                   { return true }
func DeleteAll(collection interface{}, query interface{}) bool                                   { return true }
func Find(collection interface{}, query interface{}) interface{}                                 { return nil }
func Upsert(collection interface{}, query interface{}, i interface{}) bool                       { return true }
func Get(collection interface{}, id string, i interface{})                                       {}
func Get2(collection interface{}, id interface{}, i interface{})                                 {}
func GetByQ(collection interface{}, query interface{}, i interface{})                            {}
func ListByQLimit(collection interface{}, query interface{}, i interface{}, limit int)           {}
func Distinct(collection interface{}, query map[string]interface{}, field string, i interface{}) {}
func Err(err error) bool                                                                         { return true }

func PGLog(message string) {
	fmt.Println(message)
}
