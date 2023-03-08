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

func TestSchemeDayHandler_ByScheme(t *testing.T) {
	type fields struct {
		drugRepository      *mock_database.MockDrugsRepository
		procedureRepository *mock_database.MockProceduresRepository
		schemeRepository    *mock_database.MockSchemesRepository
		schemeDayRepository *mock_database.MockSchemeDaysRepository
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
			name: "schemeDayHandler.All failed repository",
			prepare: func(f *fields) {
				f.schemeDayRepository.EXPECT().ByScheme(
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
			name: "schemeDayHandler.All returns all scheme days without pagination",
			prepare: func(f *fields) {
				f.schemeDayRepository.EXPECT().ByScheme(
					1,
					defaultPagination.Limit,
					defaultPagination.Offset(),
				).Return([]database.SchemeDay{}, nil)
			},
			error: false,
		},
		{
			name: "schemeDayHandler.All bind pagination params",
			prepare: func(f *fields) {
				f.schemeDayRepository.EXPECT().ByScheme(1, 10, 0).Return([]database.SchemeDay{}, nil)
			},
			params: "page=1&limit=10",
			error:  false,
		},
		{
			name:   "schemeDayHandler.All failed bind pagination params",
			params: "page=1&limit=test",
			error:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				drugRepository:      mock_database.NewMockDrugsRepository(ctrl),
				procedureRepository: mock_database.NewMockProceduresRepository(ctrl),
				schemeRepository:    mock_database.NewMockSchemesRepository(ctrl),
				schemeDayRepository: mock_database.NewMockSchemeDaysRepository(ctrl),
			}
			if test.prepare != nil {
				test.prepare(&f)
			}

			sdh := NewSchemeDaysHandler(
				f.paginator,
				f.drugRepository,
				f.procedureRepository,
				f.schemeRepository,
				f.schemeDayRepository,
			)

			req := httptest.NewRequest(
				"GET",
				fmt.Sprintf("/api/schemes/:schemeID/days?%s", test.params),
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()

			e := echo.New()
			NewRouter(e)

			c := e.NewContext(req, rec)
			c.SetParamNames("schemeID")
			c.SetParamValues("1")

			if err := sdh.ByScheme(c); (err != nil) != test.error {
				t.Errorf("schemeDayHandler.Create error = %v, wantErr %v", err, test.error)
			}
		})
	}
}

func TestSchemeDayHandler_Create(t *testing.T) {
	type fields struct {
		drugRepository      *mock_database.MockDrugsRepository
		procedureRepository *mock_database.MockProceduresRepository
		schemeRepository    *mock_database.MockSchemesRepository
		schemeDayRepository *mock_database.MockSchemeDaysRepository
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
			name:    "schemeDayHandler.Create invalid body",
			prepare: func(f *fields) {},
			body:    `{"schemeDay":{"test":"Test"}}`,
			error:   true,
		},
		{
			name: "schemeDayHandler.Create failed a scheme repository",
			prepare: func(f *fields) {
				f.schemeRepository.EXPECT().ByID(uint(1)).Return(database.Scheme{}, fmt.Errorf("test error"))
			},
			body:  `{"schemeDay":{"procedure":1,"drug":1,"order":1,"times":1,"frequency":1}}`,
			error: true,
		},
		{
			name: "schemeDayHandler.Create failed new scheme day when failed validation of order",
			prepare: func(f *fields) {
				f.schemeRepository.EXPECT().ByID(uint(1)).Return(database.Scheme{ID: uint(1), IllnessID: uint(1), Length: uint(1)}, nil)
			},
			body:  `{"schemeDay":{"procedure":1,"drug":1,"order":2,"times":1,"frequency":1}}`,
			error: true,
		},
		{
			name: "schemeDayHandler.Create failed new scheme day when failed validation of frequency",
			prepare: func(f *fields) {
				f.schemeRepository.EXPECT().ByID(uint(1)).Return(database.Scheme{ID: uint(1), IllnessID: uint(1), Length: uint(1)}, nil)
			},
			body:  `{"schemeDay":{"procedure":1,"drug":1,"order":1,"times":5,"frequency":6}}`,
			error: true,
		},
		{
			name: "schemeDayHandler.Create failed drug repository",
			prepare: func(f *fields) {
				f.schemeRepository.EXPECT().ByID(uint(1)).Return(database.Scheme{ID: uint(1), IllnessID: uint(1), Length: uint(1)}, nil)
				f.drugRepository.EXPECT().ByID(uint(1)).Return(database.Drug{}, fmt.Errorf("test error"))
			},
			body:  `{"schemeDay":{"procedure":1,"drug":1,"order":1,"times":1,"frequency":1}}`,
			error: true,
		},
		{
			name: "schemeDayHandler.Create failed procedure repository",
			prepare: func(f *fields) {
				f.schemeRepository.EXPECT().ByID(uint(1)).Return(database.Scheme{ID: uint(1), IllnessID: uint(1), Length: uint(1)}, nil)
				f.drugRepository.EXPECT().ByID(uint(1)).Return(database.Drug{ID: uint(1), Title: "Test"}, nil)
				f.procedureRepository.EXPECT().ByID(uint(1)).Return(database.Procedure{}, fmt.Errorf("test error"))
			},
			body:  `{"schemeDay":{"procedure":1,"drug":1,"order":1,"times":1,"frequency":1}}`,
			error: true,
		},
		{
			name: "schemeDayHandler.Create failed add method's scheme day repository",
			prepare: func(f *fields) {
				f.drugRepository.EXPECT().ByID(uint(1)).Return(database.Drug{ID: uint(1), Title: "Test"}, nil)
				f.procedureRepository.EXPECT().ByID(uint(1)).Return(database.Procedure{ID: uint(1), Title: "Test"}, nil)
				f.schemeRepository.EXPECT().ByID(uint(1)).Return(database.Scheme{ID: uint(1), IllnessID: uint(1), Length: uint(1)}, nil)
				f.schemeDayRepository.EXPECT().Add(uint(1), uint(1), uint(1), uint(1), uint(1), uint(1)).Return(uint(0), fmt.Errorf("test error"))
			},
			body:  `{"schemeDay":{"procedure":1,"drug":1,"order":1,"times":1,"frequency":1}}`,
			error: true,
		},
		{
			name: "schemeDayHandler.Create returns new scheme day",
			prepare: func(f *fields) {
				f.drugRepository.EXPECT().ByID(uint(1)).Return(database.Drug{ID: uint(1), Title: "Test"}, nil)
				f.procedureRepository.EXPECT().ByID(uint(1)).Return(database.Procedure{ID: uint(1), Title: "Test"}, nil)
				f.schemeRepository.EXPECT().ByID(uint(1)).Return(database.Scheme{ID: uint(1), IllnessID: uint(1), Length: uint(1)}, nil)
				f.schemeDayRepository.EXPECT().Add(uint(1), uint(1), uint(1), uint(1), uint(1), uint(1)).Return(uint(1), nil)
			},
			body:  `{"schemeDay":{"procedure":1,"drug":1,"order":1,"times":1,"frequency":1}}`,
			error: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				drugRepository:      mock_database.NewMockDrugsRepository(ctrl),
				procedureRepository: mock_database.NewMockProceduresRepository(ctrl),
				schemeRepository:    mock_database.NewMockSchemesRepository(ctrl),
				schemeDayRepository: mock_database.NewMockSchemeDaysRepository(ctrl),
			}
			if test.prepare != nil {
				test.prepare(&f)
			}

			sdh := NewSchemeDaysHandler(
				f.paginator,
				f.drugRepository,
				f.procedureRepository,
				f.schemeRepository,
				f.schemeDayRepository,
			)

			req := httptest.NewRequest(
				"POST",
				"/api/schemes/:schemeID/days",
				strings.NewReader(test.body),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()

			e := echo.New()
			NewRouter(e)

			c := e.NewContext(req, rec)
			c.SetParamNames("schemeID")
			c.SetParamValues("1")

			if err := sdh.CreateForScheme(c); (err != nil) != test.error {
				t.Errorf("schemeHandler.Create error = %v, wantErr %v", err, test.error)
			}
		})
	}
}
