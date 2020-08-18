package repo

import "github.com/sawickiszymon/gowebapp/models"

// Repository
type PostRepo interface {
	Create(e *models.Email) error
	ViewMessages(pageNumber int, email string) ([]models.Email, error)
	SendEmails(magicNumber int) (error, []string)
}
