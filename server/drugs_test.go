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

func TestDrugHandler_All(t *testing.T) {
	type fields struct {
		drugRepository *mock_database.MockDrugsRepository
		paginator      Pagination
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
			name: "drugHandler.All failed repository",
			prepare: func(f *fields) {
				f.drugRepository.EXPECT().All(
					defaultPagination.Limit,
					defaultPagination.Offset(),
				).Return(
					nil, fmt.Errorf("test error"),
				)
			},
			error: true,
		},
		{
			name: "drugHandler.All returns all drug without pagination",
			prepare: func(f *fields) {
				f.drugRepository.EXPECT().All(
					defaultPagination.Limit,
					defaultPagination.Offset(),
				).Return([]database.Drug{}, nil)
			},
			error: false,
		},
		{
			name: "drugHandler.All bind pagination params",
			prepare: func(f *fields) {
				f.drugRepository.EXPECT().All(10, 0).Return([]database.Drug{}, nil)
			},
			params: "page=1&limit=10",
			error:  false,
		},
		{
			name:   "drugHandler.All failed bind pagination params",
			params: "page=1&limit=test",
			error:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				drugRepository: mock_database.NewMockDrugsRepository(ctrl),
			}
			if test.prepare != nil {
				test.prepare(&f)
			}

			dh := NewDrugsHandler(f.paginator, f.drugRepository)

			req := httptest.NewRequest(
				"GET",
				fmt.Sprintf("/api/drugs?%s", test.params),
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()

			e := echo.New()
			NewRouter(e)

			c := e.NewContext(req, rec)

			if err := dh.All(c); (err != nil) != test.error {
				t.Errorf("drugHandler.Create error = %v, wantErr %v", err, test.error)
			}
		})
	}
}

func TestDrugHandler_Create(t *testing.T) {
	type fields struct {
		drugRepository *mock_database.MockDrugsRepository
		paginator      Pagination
	}

	tests := []struct {
		name    string
		prepare func(f *fields)
		body    string
		status  int
		error   bool
	}{
		{
			name:    "drugHandler.Create invalid body",
			prepare: func(f *fields) {},
			body:    `{"drug":{"test":"Test"}}`,
			error:   true,
		},
		{
			name: "drugHandler.Create failed repository",
			prepare: func(f *fields) {
				f.drugRepository.EXPECT().Add("Test").Return(uint(0), fmt.Errorf("test error"))
			},
			body:  `{"drug":{"title":"Test"}}`,
			error: true,
		},
		{
			name: "drugHandler.Create returns new drug",
			prepare: func(f *fields) {
				f.drugRepository.EXPECT().Add("Test").Return(uint(1), nil)
			},
			body:  `{"drug":{"title":"Test"}}`,
			error: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				drugRepository: mock_database.NewMockDrugsRepository(ctrl),
			}
			if test.prepare != nil {
				test.prepare(&f)
			}

			dh := NewDrugsHandler(f.paginator, f.drugRepository)

			req := httptest.NewRequest(
				"POST",
				"/api/drugs",
				strings.NewReader(test.body),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()

			e := echo.New()
			NewRouter(e)

			c := e.NewContext(req, rec)

			if err := dh.Create(c); (err != nil) != test.error {
				t.Errorf("drugHandler.Create error = %v, wantErr %v", err, test.error)
			}
		})
	}
}
