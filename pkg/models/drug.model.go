package models

type Drug struct {
	ID    uint   `gorm:"primarykey" json:"id"`
	Title string `gorm:"type:varchar(100); not null; unique_index" json:"title"`
}
