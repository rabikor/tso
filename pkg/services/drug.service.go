package services

import (
	"treatment-scheme-organizer/pkg/models"
	"treatment-scheme-organizer/pkg/repositories"
)

type DrugConfiguration func(ds *DrugService) error

type DrugService struct {
	dr repositories.DrugRepository
}

func NewDrugService(dr repositories.DrugRepository) *DrugService {
	return &DrugService{dr: dr}
}

func (ds DrugService) GetAll(limit, page uint) ([]models.Drug, error) {
	return ds.dr.GetAll(limit, page)
}

func (ds DrugService) Create(drug *models.Drug) error {
	return ds.dr.Add(drug)
}
