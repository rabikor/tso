package module

type Drug struct {
	ID    uint   `gorm:"primarykey" json:"id"`
	Title string `gorm:"type:varchar(100);not null;unique_index" json:"title"`
}

type DrugStore interface {
	GetAll(limit, page uint) ([]Drug, error)
	Add(drug *Drug) error
}
