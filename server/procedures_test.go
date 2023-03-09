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

func TestProcedureHandler_All(t *testing.T) {
	type fields struct {
		procedureRepository *mock_database.MockProceduresRepository
		paginator           Pagination
	}

	env, _ := config.NewEnv("./.env")
	defaultPagination := NewPagination(env)

	tests := []struct {
		name    string
		prepare func(f *fields)
		params  string
		status  int
		error   bool
	}{
		{
			name: "procedureHandler.All failed repository",
			prepare: func(f *fields) {
				f.procedureRepository.EXPECT().All(
					defaultPagination.Limit,
					defaultPagination.Offset(),
				).Return(
					nil, fmt.Errorf("test error"),
				)
			},
			error: true,
		},
		{
			name: "procedureHandler.All returns all procedure without pagination",
			prepare: func(f *fields) {
				f.procedureRepository.EXPECT().All(
					defaultPagination.Limit,
					defaultPagination.Offset(),
				).Return([]database.Procedure{}, nil)
			},
			error: false,
		},
		{
			name: "procedureHandler.All bind pagination params",
			prepare: func(f *fields) {
				f.procedureRepository.EXPECT().All(10, 0).Return([]database.Procedure{}, nil)
			},
			params: "page=1&limit=10",
			error:  false,
		},
		{
			name:   "procedureHandler.All failed bind pagination params",
			params: "page=1&limit=test",
			error:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				procedureRepository: mock_database.NewMockProceduresRepository(ctrl),
			}
			if test.prepare != nil {
				test.prepare(&f)
			}

			ph := NewProceduresHandler(f.paginator, f.procedureRepository)

			req := httptest.NewRequest(
				"GET",
				fmt.Sprintf("/api/procedures?%s", test.params),
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()

			e := echo.New()
			NewRouter(e)

			c := e.NewContext(req, rec)

			if err := ph.All(c); (err != nil) != test.error {
				t.Errorf("procedureHandler.Create error = %v, wantErr %v", err, test.error)
			}
		})
	}
}

func TestProcedureHandler_Create(t *testing.T) {
	type fields struct {
		procedureRepository *mock_database.MockProceduresRepository
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
			name:    "procedureHandler.Create invalid body",
			prepare: func(f *fields) {},
			body:    `{"procedure":{"test":"Test"}}`,
			error:   true,
		},
		{
			name: "procedureHandler.Create failed repository",
			prepare: func(f *fields) {
				f.procedureRepository.EXPECT().Add("Test").Return(uint(0), fmt.Errorf("test error"))
			},
			body:  `{"procedure":{"title":"Test"}}`,
			error: true,
		},
		{
			name: "procedureHandler.Create returns new procedure",
			prepare: func(f *fields) {
				f.procedureRepository.EXPECT().Add("Test").Return(uint(1), nil)
			},
			body:  `{"procedure":{"title":"Test"}}`,
			error: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				procedureRepository: mock_database.NewMockProceduresRepository(ctrl),
			}
			if test.prepare != nil {
				test.prepare(&f)
			}

			ph := NewProceduresHandler(f.paginator, f.procedureRepository)

			req := httptest.NewRequest(
				"POST",
				"/api/procedures",
				strings.NewReader(test.body),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()

			e := echo.New()
			NewRouter(e)

			c := e.NewContext(req, rec)

			if err := ph.Create(c); (err != nil) != test.error {
				t.Errorf("procedureHandler.Create error = %v, wantErr %v", err, test.error)
			}
		})
	}
}
