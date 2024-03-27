package model

type News struct {
	Id         int64   `json:"Id"`
	Title      string  `json:"Title"`
	Content    string  `json:"Content"`
	Categories []int64 `json:"Categories"`
}
