package server

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"treatment-scheme-organizer/config"
	"treatment-scheme-organizer/database"
	mock_database "treatment-scheme-organizer/mocks/database"
)

func TestSchemeHandler_ByIllness(t *testing.T) {
	type fields struct {
		illnessRepository *mock_database.MockIllnessesRepository
		schemeRepository  *mock_database.MockSchemesRepository
		paginator         Pagination
	}

	env, _ := config.NewEnv("./.env")
	defaultPagination := NewPagination(env)

	tests := []struct {
		name      string
		illnessID int
		prepare   func(f *fields)
		params    string
		status    int
		error     bool
	}{
		{
			name: "schemeHandler.All failed repository",
			prepare: func(f *fields) {
				f.schemeRepository.EXPECT().ByIllness(
					1,
					defaultPagination.Limit,
					defaultPagination.Offset(),
				).Return(
					nil, fmt.Errorf("test error"),
				)
			},
			error: true,
		},
		{
			name: "schemeHandler.All returns all scheme without pagination",
			prepare: func(f *fields) {
				f.schemeRepository.EXPECT().ByIllness(
					1,
					defaultPagination.Limit,
					defaultPagination.Offset(),
				).Return([]database.Scheme{}, nil)
			},
			error: false,
		},
		{
			name: "schemeHandler.All bind pagination params",
			prepare: func(f *fields) {
				f.schemeRepository.EXPECT().ByIllness(1, 10, 0).Return([]database.Scheme{}, nil)
			},
			params: "page=1&limit=10",
			error:  false,
		},
		{
			name:   "schemeHandler.All failed bind pagination params",
			params: "page=1&limit=test",
			error:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				illnessRepository: mock_database.NewMockIllnessesRepository(ctrl),
				schemeRepository:  mock_database.NewMockSchemesRepository(ctrl),
			}
			if test.prepare != nil {
				test.prepare(&f)
			}

			sh := NewSchemesHandler(f.paginator, f.illnessRepository, f.schemeRepository)

			req := httptest.NewRequest(
				"GET",
				fmt.Sprintf("/api/illnesses/:illnessID/schemes?%s", test.params),
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()

			e := echo.New()
			NewRouter(e)

			c := e.NewContext(req, rec)
			c.SetParamNames("illnessID")
			c.SetParamValues("1")

			if err := sh.ByIllness(c); (err != nil) != test.error {
				t.Errorf("schemeHandler.Create error = %v, wantErr %v", err, test.error)
			}
		})
	}
}

func TestSchemeHandler_Create(t *testing.T) {
	type fields struct {
		illnessRepository *mock_database.MockIllnessesRepository
		schemeRepository  *mock_database.MockSchemesRepository
		paginator         Pagination
	}

	tests := []struct {
		name    string
		prepare func(f *fields)
		body    string
		status  int
		error   bool
	}{
		{
			name:    "schemeHandler.Create invalid body",
			prepare: func(f *fields) {},
			body:    `{"scheme":{"test":"Test"}}`,
			error:   true,
		},
		{
			name: "schemeHandler.Create failed scheme repository",
			prepare: func(f *fields) {
				f.illnessRepository.EXPECT().ByID(uint(1)).Return(database.Illness{ID: uint(1), Title: "Test"}, nil)
				f.schemeRepository.EXPECT().Add(uint(1), uint(1), []database.SchemeDay{}).Return(uint(0), fmt.Errorf("test error"))
			},
			body:  `{"scheme":{"illness":1,"length":1,"days":[]}}`,
			error: true,
		},
		{
			name: "schemeHandler.Create failed illness repository",
			prepare: func(f *fields) {
				f.illnessRepository.EXPECT().ByID(uint(1)).Return(database.Illness{}, fmt.Errorf("test error"))
			},
			body:  `{"scheme":{"illness":1,"length":1,"days":[]}}`,
			error: true,
		},
		{
			name: "schemeHandler.Create returns new scheme without days",
			prepare: func(f *fields) {
				f.illnessRepository.EXPECT().ByID(uint(1)).Return(database.Illness{ID: uint(1), Title: "Test"}, nil)
				f.schemeRepository.EXPECT().Add(uint(1), uint(1), []database.SchemeDay{}).Return(uint(1), nil)
			},
			body:  `{"scheme":{"illness":1,"length":1,"days":[]}}`,
			error: false,
		},
		{
			name:  "schemeHandler.Create returns new scheme with failed order day",
			body:  `{"scheme":{"illness":1,"length":1},"days":[{"procedure":1,"drug":1,"order":2,"times":1,"frequency":1}]}`,
			error: true,
		},
		{
			name:  "schemeHandler.Create returns new scheme with failed frequency day",
			body:  `{"scheme":{"illness":1,"length":1},"days":[{"procedure":1,"drug":1,"order":1,"times":5,"frequency":6}]}`,
			error: true,
		},
		{
			name: "schemeHandler.Create returns new scheme with days",
			prepare: func(f *fields) {
				f.illnessRepository.EXPECT().ByID(uint(1)).Return(database.Illness{ID: uint(1), Title: "Test"}, nil)
				f.schemeRepository.EXPECT().Add(uint(1), uint(1), []database.SchemeDay{
					{
						ProcedureID: uint(1),
						DrugID:      uint(1),
						Order:       uint(1),
						Times:       uint(1),
						Frequency:   uint(1),
					},
				}).Return(uint(1), nil)
			},
			body:  `{"scheme":{"illness":1,"length":1},"days":[{"procedure":1,"drug":1,"order":1,"times":1,"frequency":1}]}`,
			error: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				illnessRepository: mock_database.NewMockIllnessesRepository(ctrl),
				schemeRepository:  mock_database.NewMockSchemesRepository(ctrl),
			}
			if test.prepare != nil {
				test.prepare(&f)
			}

			sh := NewSchemesHandler(f.paginator, f.illnessRepository, f.schemeRepository)

			req := httptest.NewRequest(
				"POST",
				"/api/schemes",
				strings.NewReader(test.body),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()

			e := echo.New()
			NewRouter(e)

			c := e.NewContext(req, rec)

			if err := sh.Create(c); (err != nil) != test.error {
				t.Errorf("schemeHandler.Create error = %v, wantErr %v", err, test.error)
			}
		})
	}
}
