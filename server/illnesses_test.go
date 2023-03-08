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

func TestIllnessHandler_All(t *testing.T) {
	type fields struct {
		illnessRepository *mock_database.MockIllnessesRepository
		paginator         Pagination
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
			name: "illnessHandler.All failed repository",
			prepare: func(f *fields) {
				f.illnessRepository.EXPECT().All(
					defaultPagination.Limit,
					defaultPagination.Offset(),
				).Return(
					nil, fmt.Errorf("test error"),
				)
			},
			error: true,
		},
		{
			name: "illnessHandler.All returns all illness without pagination",
			prepare: func(f *fields) {
				f.illnessRepository.EXPECT().All(
					defaultPagination.Limit,
					defaultPagination.Offset(),
				).Return([]database.Illness{}, nil)
			},
			error: false,
		},
		{
			name: "illnessHandler.All bind pagination params",
			prepare: func(f *fields) {
				f.illnessRepository.EXPECT().All(10, 0).Return([]database.Illness{}, nil)
			},
			params: "page=1&limit=10",
			error:  false,
		},
		{
			name:   "illnessHandler.All failed bind pagination params",
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
			}
			if test.prepare != nil {
				test.prepare(&f)
			}

			ih := NewIllnessesHandler(f.paginator, f.illnessRepository)

			req := httptest.NewRequest(
				"GET",
				fmt.Sprintf("/api/illnesses?%s", test.params),
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()

			e := echo.New()
			NewRouter(e)

			c := e.NewContext(req, rec)

			if err := ih.All(c); (err != nil) != test.error {
				t.Errorf("illnessHandler.Create error = %v, wantErr %v", err, test.error)
			}
		})
	}
}

func TestIllnessHandler_Create(t *testing.T) {
	type fields struct {
		illnessRepository *mock_database.MockIllnessesRepository
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
			name:    "illnessHandler.Create invalid body",
			prepare: func(f *fields) {},
			body:    `{"illness":{"test":"Test"}}`,
			error:   true,
		},
		{
			name: "illnessHandler.Create failed repository",
			prepare: func(f *fields) {
				f.illnessRepository.EXPECT().Add("Test").Return(uint(0), fmt.Errorf("test error"))
			},
			body:  `{"illness":{"title":"Test"}}`,
			error: true,
		},
		{
			name: "illnessHandler.Create returns new illness",
			prepare: func(f *fields) {
				f.illnessRepository.EXPECT().Add("Test").Return(uint(1), nil)
			},
			body:  `{"illness":{"title":"Test"}}`,
			error: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				illnessRepository: mock_database.NewMockIllnessesRepository(ctrl),
			}
			if test.prepare != nil {
				test.prepare(&f)
			}

			ih := NewIllnessesHandler(f.paginator, f.illnessRepository)

			req := httptest.NewRequest(
				"POST",
				"/api/illnesses",
				strings.NewReader(test.body),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()

			e := echo.New()
			NewRouter(e)

			c := e.NewContext(req, rec)

			if err := ih.Create(c); (err != nil) != test.error {
				t.Errorf("illnessHandler.Create error = %v, wantErr %v", err, test.error)
			}
		})
	}
}
