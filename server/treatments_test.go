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

func TestTreatmentHandler_ByIllness(t *testing.T) {
	type fields struct {
		illnessRepository   *mock_database.MockIllnessesRepository
		schemeRepository    *mock_database.MockSchemesRepository
		treatmentRepository *mock_database.MockTreatmentsRepository
		paginator           Pagination
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
			name: "treatmentHandler.All failed repository",
			prepare: func(f *fields) {
				f.treatmentRepository.EXPECT().ByIllness(
					uint(1),
					defaultPagination.Limit,
					defaultPagination.Offset(),
				).Return(
					nil, fmt.Errorf("test error"),
				)
			},
			error: true,
		},
		{
			name: "treatmentHandler.All returns all treatment without pagination",
			prepare: func(f *fields) {
				f.treatmentRepository.EXPECT().ByIllness(
					uint(1),
					defaultPagination.Limit,
					defaultPagination.Offset(),
				).Return([]database.Treatment{}, nil)
			},
			error: false,
		},
		{
			name: "treatmentHandler.All bind pagination params",
			prepare: func(f *fields) {
				f.treatmentRepository.EXPECT().ByIllness(uint(1), 10, 0).Return([]database.Treatment{}, nil)
			},
			params: "page=1&limit=10",
			error:  false,
		},
		{
			name:   "treatmentHandler.All failed bind pagination params",
			params: "page=1&limit=test",
			error:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				illnessRepository:   mock_database.NewMockIllnessesRepository(ctrl),
				schemeRepository:    mock_database.NewMockSchemesRepository(ctrl),
				treatmentRepository: mock_database.NewMockTreatmentsRepository(ctrl),
			}
			if test.prepare != nil {
				test.prepare(&f)
			}

			sh := NewTreatmentsHandler(f.paginator, f.illnessRepository, f.schemeRepository, f.treatmentRepository)

			req := httptest.NewRequest(
				"GET",
				fmt.Sprintf("/api/illnesses/:illnessID/treatments?%s", test.params),
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
				t.Errorf("treatmentHandler.Create error = %v, wantErr %v", err, test.error)
			}
		})
	}
}

func TestTreatmentHandler_Create(t *testing.T) {
	type fields struct {
		illnessRepository   *mock_database.MockIllnessesRepository
		schemeRepository    *mock_database.MockSchemesRepository
		treatmentRepository *mock_database.MockTreatmentsRepository
		paginator           Pagination
	}

	tests := []struct {
		name    string
		prepare func(f *fields)
		body    string
		status  int
		error   bool
	}{
		{
			name:    "treatmentHandler.Create invalid body",
			prepare: func(f *fields) {},
			body:    `{"treatment":{"test":"Test"}}`,
			error:   true,
		},
		{
			name: "treatmentHandler.Create failed treatment repository",
			prepare: func(f *fields) {
				f.illnessRepository.EXPECT().ByID(uint(1)).Return(database.Illness{ID: uint(1), Title: "Test"}, nil)
				f.treatmentRepository.EXPECT().Add(uint(1), "12-12-2022", "12-01-2023", []database.TreatmentScheme{}).Return(uint(0), fmt.Errorf("test error"))
			},
			body:  `{"treatment":{"illness":1,"begunAt":"12-12-2022","endedAt":"12-01-2023","schemes":[]}}`,
			error: true,
		},
		{
			name: "treatmentHandler.Create failed illness repository",
			prepare: func(f *fields) {
				f.illnessRepository.EXPECT().ByID(uint(1)).Return(database.Illness{}, fmt.Errorf("test error"))
			},
			body:  `{"treatment":{"illness":1,"begunAt":"12-12-2022","endedAt":"12-01-2023","schemes":[]}}`,
			error: true,
		},
		{
			name: "treatmentHandler.Create returns new treatment without days",
			prepare: func(f *fields) {
				f.illnessRepository.EXPECT().ByID(uint(1)).Return(database.Illness{ID: uint(1), Title: "Test"}, nil)
				f.treatmentRepository.EXPECT().Add(uint(1), "12-12-2022", "12-01-2023", []database.TreatmentScheme{}).Return(uint(1), nil)
			},
			body:  `{"treatment":{"illness":1,"begunAt":"12-12-2022","endedAt":"12-01-2023","schemes":[]}}`,
			error: false,
		},
		{
			name:  "treatmentHandler.Create returns new treatment with failed order day",
			body:  `{"treatment":{"illness":1,"begunAt":"12-12-2022","endedAt":"12-01-2023","schemes":[{"procedure":1,"drug":1,"order":2,"times":1,"frequency":1}]}`,
			error: true,
		},
		{
			name:  "treatmentHandler.Create returns new treatment with failed frequency day",
			body:  `{"treatment":{"illness":1,"begunAt":"12-12-2022","endedAt":"12-01-2023","schemes":[{"procedure":1,"drug":1,"order":1,"times":5,"frequency":6}]}`,
			error: true,
		},
		{
			name: "treatmentHandler.Create failed scheme repository in schemes",
			prepare: func(f *fields) {
				f.schemeRepository.EXPECT().ByID(uint(1)).Return(database.Scheme{}, fmt.Errorf("test error"))
			},
			body:  `{"treatment":{"illness":1,"begunAt":"12-12-2022","endedAt":"12-01-2023"},"schemes":[{"scheme":1,"beginFromDay":1}]}`,
			error: true,
		},
		{
			name: "treatmentHandler.Create failed scheme validation in schemes",
			prepare: func(f *fields) {
				f.schemeRepository.EXPECT().ByID(uint(1)).Return(database.Scheme{ID: uint(1), Length: 1}, nil)
			},
			body:  `{"treatment":{"illness":1,"begunAt":"12-12-2022","endedAt":"12-01-2023"},"schemes":[{"scheme":1,"beginFromDay":2}]}`,
			error: true,
		},
		{
			name: "treatmentHandler.Create returns new treatment with days",
			prepare: func(f *fields) {
				f.illnessRepository.EXPECT().ByID(uint(1)).Return(database.Illness{ID: uint(1), Title: "Test"}, nil)
				f.schemeRepository.EXPECT().ByID(uint(1)).Return(database.Scheme{ID: uint(1), Length: 1}, nil)
				f.treatmentRepository.EXPECT().Add(uint(1), "12-12-2022", "12-01-2023", []database.TreatmentScheme{
					{
						SchemeID:     uint(1),
						BeginFromDay: uint(1),
						Order:        uint(1),
					},
				}).Return(uint(1), nil)
			},
			body:  `{"treatment":{"illness":1,"begunAt":"12-12-2022","endedAt":"12-01-2023"},"schemes":[{"scheme":1,"beginFromDay":1}]}`,
			error: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				illnessRepository:   mock_database.NewMockIllnessesRepository(ctrl),
				schemeRepository:    mock_database.NewMockSchemesRepository(ctrl),
				treatmentRepository: mock_database.NewMockTreatmentsRepository(ctrl),
			}
			if test.prepare != nil {
				test.prepare(&f)
			}

			sh := NewTreatmentsHandler(f.paginator, f.illnessRepository, f.schemeRepository, f.treatmentRepository)

			req := httptest.NewRequest(
				"POST",
				"/api/treatments",
				strings.NewReader(test.body),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()

			e := echo.New()
			NewRouter(e)

			c := e.NewContext(req, rec)

			if err := sh.Create(c); (err != nil) != test.error {
				t.Errorf("treatmentHandler.Create error = %v, wantErr %v", err, test.error)
			}
		})
	}
}
