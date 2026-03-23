package mongodb

import (
	"gopkg.in/mgo.v2"

	"github.com/leanote/leanote/app/db"
)

type MongoTransaction struct {
	session *mgo.Session
}

func NewMongoTransaction(session *mgo.Session) db.Transaction {
	return &MongoTransaction{session: session}
}

func (mt *MongoTransaction) Commit() error {
	return nil
}

func (mt *MongoTransaction) Rollback() error {
	return nil
}

func (mt *MongoTransaction) Exec(query string, args ...interface{}) (db.Result, error) {
	return nil, nil
}

func (mt *MongoTransaction) Query(query string, args ...interface{}) (db.Rows, error) {
	return nil, nil
}

func (mt *MongoTransaction) QueryRow(query string, args ...interface{}) db.Row {
	return nil
}
