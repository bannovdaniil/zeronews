package model

//go:generate reform

// reform:news
type News struct {
	Id         int64   `json:"Id" reform:"id,pk"`
	Title      string  `json:"Title" reform:"title"`
	Content    string  `json:"Content" reform:"content"`
	Categories []int64 `json:"Categories" reform:"-"`
}
