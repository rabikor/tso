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

func TestTreatmentSchemeHandler_ByTreatment(t *testing.T) {
	type fields struct {
		schemeRepository          *mock_database.MockSchemesRepository
		treatmentRepository       *mock_database.MockTreatmentsRepository
		treatmentSchemeRepository *mock_database.MockTreatmentSchemesRepository
		paginator                 Pagination
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
			name: "treatmentSchemeHandler.All failed repository",
			prepare: func(f *fields) {
				f.treatmentSchemeRepository.EXPECT().ByTreatment(
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
			name: "treatmentSchemeHandler.All returns all treatment schemes without pagination",
			prepare: func(f *fields) {
				f.treatmentSchemeRepository.EXPECT().ByTreatment(
					uint(1),
					defaultPagination.Limit,
					defaultPagination.Offset(),
				).Return([]database.TreatmentScheme{}, nil)
			},
			error: false,
		},
		{
			name: "treatmentSchemeHandler.All bind pagination params",
			prepare: func(f *fields) {
				f.treatmentSchemeRepository.EXPECT().ByTreatment(uint(1), 10, 0).Return([]database.TreatmentScheme{}, nil)
			},
			params: "page=1&limit=10",
			error:  false,
		},
		{
			name:   "treatmentSchemeHandler.All failed bind pagination params",
			params: "page=1&limit=test",
			error:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				schemeRepository:          mock_database.NewMockSchemesRepository(ctrl),
				treatmentRepository:       mock_database.NewMockTreatmentsRepository(ctrl),
				treatmentSchemeRepository: mock_database.NewMockTreatmentSchemesRepository(ctrl),
			}
			if test.prepare != nil {
				test.prepare(&f)
			}

			tsh := NewTreatmentSchemesHandler(
				f.paginator,
				f.schemeRepository,
				f.treatmentRepository,
				f.treatmentSchemeRepository,
			)

			req := httptest.NewRequest(
				"GET",
				fmt.Sprintf("/api/treatments/:treatmentID/schemes?%s", test.params),
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()

			e := echo.New()
			NewRouter(e)

			c := e.NewContext(req, rec)
			c.SetParamNames("treatmentID")
			c.SetParamValues("1")

			if err := tsh.ByTreatment(c); (err != nil) != test.error {
				t.Errorf("treatmentSchemeHandler.Create error = %v, wantErr %v", err, test.error)
			}
		})
	}
}

func TestTreatmentSchemeHandler_Create(t *testing.T) {
	type fields struct {
		schemeRepository          *mock_database.MockSchemesRepository
		treatmentRepository       *mock_database.MockTreatmentsRepository
		treatmentSchemeRepository *mock_database.MockTreatmentSchemesRepository
		paginator                 Pagination
	}

	tests := []struct {
		name    string
		prepare func(f *fields)
		body    string
		status  int
		error   bool
	}{
		{
			name:    "treatmentSchemeHandler.Create invalid body",
			prepare: func(f *fields) {},
			body:    `{"treatmentScheme":{"test":"Test"}}`,
			error:   true,
		},
		{
			name: "treatmentSchemeHandler.Create failed a treatment repository",
			prepare: func(f *fields) {
				f.treatmentRepository.EXPECT().ByID(uint(1)).Return(database.Treatment{}, fmt.Errorf("test error"))
			},
			body:  `{"treatmentScheme":{"scheme":1,"beginFromDay":1}}`,
			error: true,
		},
		{
			name: "treatmentSchemeHandler.Create failed a scheme repository",
			prepare: func(f *fields) {
				f.treatmentRepository.EXPECT().ByID(uint(1)).Return(database.Treatment{}, nil)
				f.schemeRepository.EXPECT().ByID(uint(1)).Return(database.Scheme{}, fmt.Errorf("test error"))
			},
			body:  `{"treatmentScheme":{"scheme":1,"beginFromDay":1}}`,
			error: true,
		},
		{
			name: "treatmentSchemeHandler.Create failed new treatment scheme when begin from not exists day",
			prepare: func(f *fields) {
				f.treatmentRepository.EXPECT().ByID(uint(1)).Return(database.Treatment{ID: uint(1), IllnessID: uint(1), BegunAt: "12-12-2022", EndedAt: "12-01-2023"}, nil)
				f.schemeRepository.EXPECT().ByID(uint(1)).Return(database.Scheme{ID: uint(1), IllnessID: uint(1), Length: uint(1)}, nil)
			},
			body:  `{"treatmentScheme":{"scheme":1,"beginFromDay":2}}`,
			error: true,
		},
		{
			name: "treatmentSchemeHandler.Create failed getting count treatment schemes by treatment",
			prepare: func(f *fields) {
				f.treatmentRepository.EXPECT().ByID(uint(1)).Return(database.Treatment{ID: uint(1), IllnessID: uint(1), BegunAt: "12-12-2022", EndedAt: "12-01-2023"}, nil)
				f.schemeRepository.EXPECT().ByID(uint(1)).Return(database.Scheme{ID: uint(1), IllnessID: uint(1), Length: uint(1)}, nil)
				f.treatmentSchemeRepository.EXPECT().CountByTreatment(uint(1)).Return(uint(0), fmt.Errorf("test error"))
			},
			body:  `{"treatmentScheme":{"scheme":1,"beginFromDay":1}}`,
			error: true,
		},
		{
			name: "treatmentSchemeHandler.Create failed add method's treatment scheme repository",
			prepare: func(f *fields) {
				f.treatmentRepository.EXPECT().ByID(uint(1)).Return(database.Treatment{ID: uint(1), IllnessID: uint(1), BegunAt: "12-12-2022", EndedAt: "12-01-2023"}, nil)
				f.schemeRepository.EXPECT().ByID(uint(1)).Return(database.Scheme{ID: uint(1), IllnessID: uint(1), Length: uint(1)}, nil)
				f.treatmentSchemeRepository.EXPECT().CountByTreatment(uint(1)).Return(uint(0), nil)
				f.treatmentSchemeRepository.EXPECT().Add(uint(1), uint(1), uint(1), uint(1)).Return(uint(0), fmt.Errorf("test error"))
			},
			body:  `{"treatmentScheme":{"scheme":1,"beginFromDay":1}}`,
			error: true,
		},
		{
			name: "treatmentSchemeHandler.Create returns new FIRST treatment scheme in a treatment",
			prepare: func(f *fields) {
				f.treatmentRepository.EXPECT().ByID(uint(1)).Return(database.Treatment{ID: uint(1), IllnessID: uint(1), BegunAt: "12-12-2022", EndedAt: "12-01-2023"}, nil)
				f.schemeRepository.EXPECT().ByID(uint(1)).Return(database.Scheme{ID: uint(1), IllnessID: uint(1), Length: uint(1)}, nil)
				f.treatmentSchemeRepository.EXPECT().CountByTreatment(uint(1)).Return(uint(0), nil)
				f.treatmentSchemeRepository.EXPECT().Add(uint(1), uint(1), uint(1), uint(1)).Return(uint(1), nil)
			},
			body:  `{"treatmentScheme":{"scheme":1,"beginFromDay":1}}`,
			error: false,
		},
		{
			name: "treatmentSchemeHandler.Create returns new SECOND treatment scheme",
			prepare: func(f *fields) {
				f.treatmentRepository.EXPECT().ByID(uint(1)).Return(database.Treatment{ID: uint(1), IllnessID: uint(1), BegunAt: "12-12-2022", EndedAt: "12-01-2023"}, nil)
				f.schemeRepository.EXPECT().ByID(uint(1)).Return(database.Scheme{ID: uint(1), IllnessID: uint(1), Length: uint(1)}, nil)
				f.treatmentSchemeRepository.EXPECT().CountByTreatment(uint(1)).Return(uint(1), nil)
				f.treatmentSchemeRepository.EXPECT().Add(uint(1), uint(1), uint(1), uint(2)).Return(uint(1), nil)
			},
			body:  `{"treatmentScheme":{"scheme":1,"beginFromDay":1}}`,
			error: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				schemeRepository:          mock_database.NewMockSchemesRepository(ctrl),
				treatmentRepository:       mock_database.NewMockTreatmentsRepository(ctrl),
				treatmentSchemeRepository: mock_database.NewMockTreatmentSchemesRepository(ctrl),
			}
			if test.prepare != nil {
				test.prepare(&f)
			}

			tsh := NewTreatmentSchemesHandler(
				f.paginator,
				f.schemeRepository,
				f.treatmentRepository,
				f.treatmentSchemeRepository,
			)

			req := httptest.NewRequest(
				"POST",
				"/api/treatments/:treatmentID/schemes",
				strings.NewReader(test.body),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()

			e := echo.New()
			NewRouter(e)

			c := e.NewContext(req, rec)
			c.SetParamNames("treatmentID")
			c.SetParamValues("1")

			if err := tsh.CreateForTreatment(c); (err != nil) != test.error {
				t.Errorf("schemeHandler.Create error = %v, wantErr %v", err, test.error)
			}
		})
	}
}
