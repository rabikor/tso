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

func createIllnessHandler() IllnessHandler {
	env := &config.Env{}
	_ = env.ParseEnv("../.env")

	db, _ := database.Open(env)

	return NewIllnessesHandler(env, db)
}

func TestNewIllnessesHandler(t *testing.T) {
	ih := createIllnessHandler()

	assert.IsType(t, IllnessHandler{}, ih)
}

func TestIllnessHandler_GetAll(t *testing.T) {
	tearDown()
	setup()

	ih := createIllnessHandler()

	e := router.New()

	req := httptest.NewRequest("GET", "/api/illnesses", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	assert.NoError(t, ih.GetAll(c))
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

func TestIllnessHandler_Create(t *testing.T) {
	tearDown()
	setup()

	ih := createIllnessHandler()

	e := router.New()
	req := httptest.NewRequest(
		"POST",
		"/api/illnesses",
		strings.NewReader(`{"illness":{"title": "Illness 4"}}`),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	assert.NoError(t, ih.Create(c))
	if assert.Equal(t, http.StatusCreated, rec.Code) {
		iList, err := d.Illnesses.GetAll(10, 0)
		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, 4, len(iList))
	}
}
