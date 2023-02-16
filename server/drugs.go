package server

import (
	"net/http"
	"strconv"

	"treatment-scheme-organizer/config"
	"treatment-scheme-organizer/database"

	"github.com/gin-gonic/gin"
)

type DrugHandler struct {
	db *database.DB
}

func NewDrugsHandler(db *database.DB) DrugHandler {
	return DrugHandler{db: db}
}

func (h DrugHandler) AddRoutes(rg *gin.RouterGroup) {
	router := rg.Group("drugs")
	router.GET("", h.GetAll)
	router.POST("", h.Create)
}

func (h DrugHandler) GetAll(c *gin.Context) {
	var (
		limit  = config.Env.API.Request.Limit
		offset = config.Env.API.Request.Offset
	)

	if x, ok := c.GetQuery("limit"); ok {
		limit, _ = strconv.Atoi(x)
	}
	if x, ok := c.GetQuery("offset"); ok {
		offset, _ = strconv.Atoi(x)
	}

	drugs, err := h.db.Drugs.GetAll(limit, offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true, "data": drugs})
}

type createDrugRequest struct {
	Title string `json:"title" binding:"required"`
}

func (h DrugHandler) Create(c *gin.Context) {
	var req createDrugRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "slug": "drug.create.bind-json", "error": err})
		return
	}

	drug := database.Drug{Title: req.Title}
	if err := h.db.Create(&drug); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "slug": "drug.create.service-request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true, "data": drug})
}
