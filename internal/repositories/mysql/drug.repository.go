package mysql

import (
	"gorm.io/gorm"
	"treatment-scheme-organizer/pkg/models"
)

type DrugRepository struct {
	db *gorm.DB
}

func NewDrugRepository(DB *gorm.DB) *DrugRepository {
	return &DrugRepository{db: DB}
}

func (dr *DrugRepository) GetAll(limit, page uint) ([]models.Drug, error) {
	var drugs []models.Drug

	result := dr.db.Limit(int(limit)).Offset(int((page - 1) * limit)).Find(&drugs)

	if result.Error != nil {
		panic(result.Error)
	}

	return drugs, nil
}

func (dr *DrugRepository) Add(drug *models.Drug) error {
	result := dr.db.Create(&drug)

	if result.Error != nil {
		panic(result.Error)
	}

	return nil
}
