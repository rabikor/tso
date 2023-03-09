package server

import "treatment-scheme-organizer/config"

type Pagination struct {
	Limit int `query:"limit" param:"limit" json:"limit"`
	Page  int `query:"page" param:"page" json:"page"`
}

func (p Pagination) Offset() int {
	return (p.Page - 1) * p.Limit
}

func NewPagination(env config.Env) Pagination {
	return Pagination{Limit: env.API.Request.Limit, Page: env.API.Request.Page}
}
