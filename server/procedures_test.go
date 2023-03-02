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

func createProcedureHandler() ProcedureHandler {
	env := &config.Env{}
	_ = env.ParseEnv("../.env")

	db, _ := database.Open(env)

	return NewProceduresHandler(env, db)
}

func TestNewProceduresHandler(t *testing.T) {
	ph := createProcedureHandler()

	assert.IsType(t, ProcedureHandler{}, ph)
}

func TestProcedureHandler_GetAll(t *testing.T) {
	tearDown()
	setup()

	ph := createProcedureHandler()

	e := router.New()

	req := httptest.NewRequest("GET", "/api/procedures", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	assert.NoError(t, ph.GetAll(c))
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

func TestProcedureHandler_Create(t *testing.T) {
	tearDown()
	setup()

	ph := createProcedureHandler()

	e := router.New()
	req := httptest.NewRequest(
		"POST",
		"/api/procedures",
		strings.NewReader(`{"procedure":{"title": "Procedure 3"}}`),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	assert.NoError(t, ph.Create(c))
	if assert.Equal(t, http.StatusCreated, rec.Code) {
		pList, err := d.Procedures.GetAll(3, 0)
		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, 3, len(pList))
	}
}
