package info

import (
	"time"
)

type EachSharedNote struct {
	Seq int
}
type EachSharedNotebook struct {
	Seq int
}

type EachSharedNotebookAndNotes struct {
	Seq             int
	DefaultNotebook map[string]EachSharedNote
	Notebooks       map[string]EachSharedNotebook
}

type SharedNotebookAndNotes struct {
	UserId string
	Shared map[string]EachSharedNotebookAndNotes
}

type SharingNotebookAndNotes struct {
	UserId    string
	Notes     map[string][]string
	Notebooks map[string][]string
}

type ShareNotebook struct {
	ShareNotebookId string    `db:"id"`
	UserId          string    `db:"user_id"`
	ToUserId        string    `db:"to_user_id"`
	NotebookId      string    `db:"notebook_id"`
	Seq             int       `db:"seq"`
	Perm            int       `db:"permissions"`
	CreatedTime     time.Time `db:"created_time"`
	IsDeleted       bool      `db:"is_deleted"`
}

type SubShareNotebooks []ShareNotebooks
type ShareNotebooks struct {
	Notebook
	ShareNotebook
	Subs SubShareNotebooks

	Seq        int
	NotebookId string
	IsDefault  bool
}

func (this SubShareNotebooks) Len() int {
	return len(this)
}
func (this SubShareNotebooks) Less(i, j int) bool {
	return this[i].ShareNotebook.Seq < this[j].ShareNotebook.Seq
}
func (this SubShareNotebooks) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

type ShareNotebooksByUser map[string][]ShareNotebooks

type ShareNote struct {
	ShareNoteId string    `db:"id"`
	UserId      string    `db:"user_id"`
	ToUserId    string    `db:"to_user_id"`
	NoteId      string    `db:"note_id"`
	Perm        int       `db:"permissions"`
	CreatedTime time.Time `db:"created_time"`
	IsDeleted   bool      `db:"is_deleted"`
}

type HasShareNote struct {
	HasShareNotebookId string `db:"id"`
	UserId             string `db:"user_id"`
	ToUserId           string `db:"to_user_id"`
	Seq                int    `db:"seq"`
}

type ShareNoteWithPerm struct {
	Note
	Perm int
}

type ShareUserInfo struct {
	ToUserId          string
	Email             string
	Perm              int
	NotebookHasShared bool
}
