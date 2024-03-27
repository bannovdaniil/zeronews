package model

type Category struct {
	Id         int64 `json:"Id"`
	NewsId     int64 `json:"-"`
	CategoryId int64 `json:"Id"`
}
