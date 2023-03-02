package server

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"treatment-scheme-organizer/config"
	"treatment-scheme-organizer/database"
	"treatment-scheme-organizer/router"
)

func createSchemeDayHandler() SchemeDayHandler {
	env := &config.Env{}
	_ = env.ParseEnv("../.env")

	db, _ := database.Open(env)

	return NewSchemeDaysHandler(env, db)
}

func TestNewSchemeDaysHandler(t *testing.T) {
	sdh := createSchemeDayHandler()

	assert.IsType(t, SchemeDayHandler{}, sdh)
}

func TestSchemeDayHandler_GetByExistsScheme(t *testing.T) {
	tearDown()
	setup()

	sdh := createSchemeDayHandler()

	e := router.New()

	req := httptest.NewRequest("GET", "/api/schemes/:schemeID/days", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/schemes/:schemeID/days")
	c.SetParamNames("schemeID")
	c.SetParamValues(strconv.FormatUint(uint64(SchemeExample.ID), 10))

	assert.NoError(t, sdh.GetByScheme(c))
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

func TestSchemeDayHandler_GetByExistsSchemeWithoutDays(t *testing.T) {
	tearDown()
	setup()

	sdh := createSchemeDayHandler()

	e := router.New()

	req := httptest.NewRequest("GET", "/api/schemes/:schemeID/days", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/schemes/:schemeID/days")
	c.SetParamNames("schemeID")
	c.SetParamValues("2")

	assert.NoError(t, sdh.GetByScheme(c))
	if assert.Equal(t, http.StatusOK, rec.Code) {
		var r struct {
			Data []interface{}
			Meta map[string]interface{}
		}

		err := json.Unmarshal(rec.Body.Bytes(), &r)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(r.Data))
	}
}

func TestSchemeDayHandler_GetByNotExistsScheme(t *testing.T) {
	tearDown()
	setup()

	sdh := createSchemeDayHandler()

	e := router.New()

	req := httptest.NewRequest("GET", "/api/schemes/:schemeID/days", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/schemes/:schemeID/days")
	c.SetParamNames("schemeID")
	c.SetParamValues("1000")

	assert.NoError(t, sdh.GetByScheme(c))
	if assert.Equal(t, http.StatusOK, rec.Code) {
		var r struct {
			Data []interface{}
			Meta map[string]interface{}
		}

		err := json.Unmarshal(rec.Body.Bytes(), &r)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(r.Data))
	}
}

func TestSchemeDayHandler_CreateToExistsScheme(t *testing.T) {
	tearDown()
	setup()

	sdh := createSchemeDayHandler()

	e := router.New()
	req := httptest.NewRequest(
		"POST",
		"/api/schemes/:schemeID/days",
		strings.NewReader(`{"schemeDay":{"drugId": 1, "procedureId": 1, "order": 3, "times": 1, "frequency": 1}}`),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/schemes/:schemeID/days")
	c.SetParamNames("schemeID")
	c.SetParamValues(strconv.FormatUint(uint64(SchemeExample.ID), 10))

	assert.NoError(t, sdh.CreateForScheme(c))
	if assert.Equal(t, http.StatusCreated, rec.Code) {
		pList, err := d.SchemeDays.GetByScheme(int(SchemeExample.ID), 3, 0)
		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, 3, len(pList))
	}
}

func TestSchemeDayHandler_CreateToNotExistsScheme(t *testing.T) {
	tearDown()
	setup()

	sdh := createSchemeDayHandler()

	e := router.New()
	req := httptest.NewRequest(
		"POST",
		"/api/schemes/:schemeID/days",
		strings.NewReader(`{"schemeDay":{"drugId": 1, "procedureId": 1, "order": 3, "times": 1, "frequency": 1}}`),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/schemes/:schemeID/days")
	c.SetParamNames("schemeID")
	c.SetParamValues("1000")

	assert.Error(t, sdh.CreateForScheme(c))
}
