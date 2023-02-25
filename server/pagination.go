package server

type Pagination struct {
	Limit int `query:"limit" param:"limit" json:"limit"`
	Page  int `query:"page" param:"page" json:"page"`
}

func (p Pagination) GetOffset() int {
	return (p.Page - 1) * p.Limit
}
