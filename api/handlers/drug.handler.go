package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"treatment-scheme-organizer/internal/configs"
	"treatment-scheme-organizer/pkg/models"
	"treatment-scheme-organizer/pkg/services"
)

type DrugHandler struct {
	ds *services.DrugService
}

type createDrugRequest struct {
	Title string `json:"title" binding:"required"`
}

func NewDrugHandler(ds *services.DrugService) DrugHandler {
	return DrugHandler{ds: ds}
}

func (dh DrugHandler) GetAll(c *gin.Context) {
	var page, perPage = configs.Env.API.Request.Page, configs.Env.API.Request.PerPage

	if pageParam, ok := c.GetQuery("page"); !ok && pageParam != "" {
		var (
			err        error
			pageUint64 uint64
		)

		if pageUint64, err = strconv.ParseUint(pageParam, 10, 8); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": err.Error()})
			return
		}

		page = uint(pageUint64)
	}

	if perPageParam, ok := c.GetQuery("perPage"); ok && perPageParam != "" {
		var (
			err           error
			perPageUint64 uint64
		)

		if perPageUint64, err = strconv.ParseUint(perPageParam, 10, 8); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": err.Error()})
			return
		}

		page = uint(perPageUint64)
	}

	drugs, err := dh.ds.GetAll(perPage, page)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true, "data": drugs})
}

func (dh DrugHandler) Create(c *gin.Context) {
	var input createDrugRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "slug": "drug.create.bind-json", "error": err.Error()})
		return
	}

	drug := models.Drug{Title: input.Title}

	if err := dh.ds.Create(&drug); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "slug": "drug.create.service-request", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true, "data": drug})
}
