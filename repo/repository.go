package repo

import "github.com/sawickiszymon/gowebapp/models"

type PostRepo interface {
	Create(e *models.Email) error
	View(pageNumber int, email string) ([]models.Email, error)
}
