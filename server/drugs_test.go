package server

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"treatment-scheme-organizer/config"
	"treatment-scheme-organizer/database"
	"treatment-scheme-organizer/router"
)

func createDrugHandler() DrugHandler {
	env := &config.Env{}
	_ = env.ParseEnv("../.env")

	db, _ := database.Open(env)

	return NewDrugsHandler(env, db)
}

func TestNewDrugsHandler(t *testing.T) {
	dh := createDrugHandler()

	assert.IsType(t, DrugHandler{}, dh)
}

func TestDrugHandler_AddRoutes(t *testing.T) {
	dh := createDrugHandler()

	e := echo.New()

	g := e.Group("/")

	dh.AddRoutes(g)

	assert.True(t, true)
}

func TestDrugHandler_GetAll(t *testing.T) {
	tearDown()
	setup()

	dh := createDrugHandler()

	e := router.New()

	req := httptest.NewRequest("GET", "/api/drugs", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	assert.NoError(t, dh.GetAll(c))
	if assert.Equal(t, http.StatusOK, rec.Code) {
		var r struct {
			Data []interface{}
			Meta map[string]interface{}
		}

		err := json.Unmarshal(rec.Body.Bytes(), &r)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(r.Data))
	}
}

func TestDrugHandler_Create(t *testing.T) {
	tearDown()
	setup()

	dh := createDrugHandler()

	e := router.New()
	req := httptest.NewRequest(
		"POST",
		"/api/drugs",
		strings.NewReader(`{"drug":{"title": "Drug 3"}}`),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	assert.NoError(t, dh.Create(c))
	if assert.Equal(t, http.StatusCreated, rec.Code) {
		dr, err := d.Drugs.GetAll(3, 0)
		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, 3, len(dr))
	}
}
