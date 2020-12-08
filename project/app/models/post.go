package models

type Post struct {
	Id            int64   `json:"id,omitempty"`
	Parent        int64   `json:"parent"`
	Thread        int32   `json:"thread,omitempty"`
	Forum         string  `json:"forum,omitempty"`
	Author        string  `json:"author"`
	Created       string  `json:"created,omitempty"`
	IsEdited      bool    `json:"isEdited"`
	Message       string  `json:"message"`
	Path          []int64 `json:"-"`
	//Childs        Posts   `json:"childs,omitempty"`
	//ParentPointer *Post   `json:"-"`
}

type PostFull struct {
	Author *User   `json:"author,omitempty"`
	Forum  *Forum  `json:"forum,omitempty"`
	Post   *Post   `json:"post,omitempty"`
	Thread *Thread `json:"thread,omitempty"`
}

type PostRelatedFields struct {
	WithUser bool
	WithForum bool
	WithThread bool
}

type Posts []*Post
