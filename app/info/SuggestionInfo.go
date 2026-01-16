package info

type Suggestion struct {
	Id         string `db:"id"`
	UserId     string `db:"user_id"`
	Content    string `db:"content"`
}
