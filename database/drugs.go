package database

import "gorm.io/gorm"

type (
	drugsRepository interface {
		GetAll(limit, page int) ([]Drug, error)
	}
	drugsTable struct {
		*gorm.DB
	}
)

type Drug struct {
	ID    uint   `gorm:"primarykey" json:"id"`
	Title string `gorm:"type:varchar(100); not null; unique_index" json:"title"`
}

func (db *drugsTable) GetAll(limit, offset int) (ds []Drug, _ error) {
	return ds, db.Find(&ds).Limit(limit).Offset(offset).Error
}
