package models

type CreateNotice struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Type    string `json:"type"`
}
