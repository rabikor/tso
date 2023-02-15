package repositories

import "treatment-scheme-organizer/pkg/models"

type DrugRepository interface {
	GetAll(limit, page uint) ([]models.Drug, error)
	Add(drug *models.Drug) error
}
