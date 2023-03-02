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

func createSchemeHandler() SchemeHandler {
	env := &config.Env{}
	_ = env.ParseEnv("../.env")

	db, _ := database.Open(env)

	return NewSchemesHandler(env, db)
}

func TestNewSchemesHandler(t *testing.T) {
	sh := createSchemeHandler()

	assert.IsType(t, SchemeHandler{}, sh)
}

func TestSchemeHandler_GetByIllness(t *testing.T) {
	tearDown()
	setup()

	sh := createSchemeHandler()

	e := router.New()

	req := httptest.NewRequest("GET", "/api/illnesses/:illnessID/schemes", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/illnesses/:illnessID/schemes")
	c.SetParamNames("illnessID")
	c.SetParamValues("1")

	assert.NoError(t, sh.GetByIllness(c))
	if assert.Equal(t, http.StatusOK, rec.Code) {
		var r struct {
			Data []interface{}
			Meta map[string]interface{}
		}

		err := json.Unmarshal(rec.Body.Bytes(), &r)
		assert.NoError(t, err)
		assert.Equal(t, 3, len(r.Data))
	}
}

func TestSchemeHandler_Create(t *testing.T) {
	tearDown()
	setup()

	sh := createSchemeHandler()

	e := router.New()
	req := httptest.NewRequest(
		"POST",
		"/api/illnesses/:illness/schemes",
		strings.NewReader(`{"scheme":{"illness": 1, "length": 1}}`),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/illnesses/:illness/schemes")
	c.SetParamNames("illness")
	c.SetParamValues("1")

	assert.NoError(t, sh.Create(c))
	if assert.Equal(t, http.StatusCreated, rec.Code) {
		sc, err := d.Schemes.GetByIllness(1, 10, 0)
		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, 4, len(sc))
	}
}
