package model

type Category struct {
	NewsId     int64 `json:"-"`
	CategoryId int64 `json:"Id"`
}
