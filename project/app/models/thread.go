package models

type Thread struct {
	Id      int32  `json:"id,omitempty"`
	Forum   string `json:"forum"`
	Author  string `json:"author"`
	Created string `json:"created"`
	Message string `json:"message"`
	Title   string `json:"title"`
	Votes   int64  `json:"votes"`
	Slug    string `json:"slug"`
}

type Threads []Thread
