package models

// Email struct - models database table
type Email struct {
	Email       string `json:"email"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	MagicNumber int    `json:"magic_number"`
}
