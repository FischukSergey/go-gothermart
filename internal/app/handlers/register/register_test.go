package register

import (
	"bytes"
	"context"
	"github.com/FischukSergey/go-gothermart.git/internal/app/handlers/mock"
	"github.com/FischukSergey/go-gothermart.git/internal/models"

	//"github.com/FischukSergey/go-gothermart.git/internal/app/middleware/auth"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestRegister(t *testing.T) {

	var log = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	type want struct {
		statusCode   int
		userID       int
		bodeResponse string
		mockError    error
	}

	tests := []struct {
		name        string
		bodyRequest string
		want        want
	}{
		{
			name:        "register successfully",
			bodyRequest: `{"login": "example","password": "password"}`,
			want: want{
				statusCode: 200,
				userID:     1,
			},
		},
		{
			name:        "login empty",
			bodyRequest: `{"login": "","password": "password"}`,
			want: want{
				statusCode:   400,
				userID:       1,
				bodeResponse: "",
			},
		},
		{
			name:        "password empty",
			bodyRequest: `{"login": "example","password": ""}`,
			want: want{
				statusCode: 400,
				userID:     1,
			},
		},
		{
			name:        "body request empty",
			bodyRequest: ``,
			want: want{
				statusCode: 400,
				userID:     1,
			},
		},
		{
			name:        "login already exists",
			bodyRequest: `{"login": "example","password": "password"}`,
			want: want{
				statusCode: 409,
				userID:     0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			s := mock.NewMockUserRegister(ctrl) //новый storage
			defer ctrl.Finish()

			requestTest := httptest.NewRequest(http.MethodPost,
				"/api/user/register",
				bytes.NewReader([]byte(tt.bodyRequest)))
			request := requestTest.WithContext(context.Background())
			w := httptest.NewRecorder()
			request.Header.Set("Content-Type", "application/json")

			switch {
			case tt.want.mockError != nil:
				s.EXPECT()
			case tt.name == "register successfully":
				s.EXPECT().Register(gomock.Any(), gomock.Any()).Return(tt.want.userID, tt.want.mockError)
			case tt.name == "login already exists":
				s.EXPECT().Register(gomock.Any(), gomock.Any()).Return(0, models.ErrUserExists)
			default:
				s.EXPECT()
			}

			h := http.HandlerFunc(Register(log, s))
			h(w, request)

			result := w.Result()
			body, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			//assert.Equal(t, tt.want.userID, 1)
			_ = body
			//assert.Contains(t, string(body), tt.want.bodyResponse)
		})
	}
}
