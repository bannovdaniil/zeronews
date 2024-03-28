package model

//go:generate reform

//reform:news_categories
type Category struct {
	Id         int64 `reform:"id,pk"`
	NewsId     int64 `reform:"news_id"`
	CategoryId int64 `reform:"category_id"`
}
