package repo

import "gowebapp/models"

type PostRepo interface {
	Create(e *models.Email) error
}
