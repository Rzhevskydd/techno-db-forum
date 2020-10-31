package models

type User struct {
	Id       int64  `json:"-"`
	Nickname string `json:"nickname"`
	FullName string `json:"fullname"`
	Email    string `json:"email"`
	About    string `json:"about"`
}

type Users []User
