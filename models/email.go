package models

type Email struct {
	Email       string `json:"email"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	MagicNumber int    `json:"magic_number"`
}
