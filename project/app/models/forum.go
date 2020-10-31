package models

type Forum struct {
	Id      int64  `json:"-"`
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	User    string `json:"user"`
	Threads int64  `json:"threads"`
	Posts   int64  `json:"posts"`
}