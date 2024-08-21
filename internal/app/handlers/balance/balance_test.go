package balance

import (
	"context"
	"errors"
	"github.com/FischukSergey/go-gothermart.git/internal/app/handlers/mock"
	"github.com/FischukSergey/go-gothermart.git/internal/app/middleware/auth"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetBalance(t *testing.T) {
	var log = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	type want struct {
		statusCode   int
		bodyResponse string
		current      float32
		withdrawn    float32
		mockError    error
	}

	tests := []struct {
		name   string
		userID int
		want   want
	}{
		{
			name:   "balance successfully",
			userID: 55,
			want: want{
				statusCode:   200,
				bodyResponse: `{"current":500.5,"withdrawn":42}` + "\n",
				current:      500.5,
				withdrawn:    42,
			},
		},
		{
			name:   "internal server error",
			userID: 55,
			want: want{
				statusCode:   500,
				bodyResponse: `{"error":"internal server error"}` + "\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			s := mock.NewMockBalancer(ctrl) //новый storage
			defer ctrl.Finish()

			requestTest := httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", nil)
			request := requestTest.WithContext(context.WithValue(requestTest.Context(),
				auth.CtxKeyUser, tt.userID))
			w := httptest.NewRecorder()
			request.Header.Set("Content-Type", "application/json")

			switch {
			case tt.want.mockError != nil:
				s.EXPECT()
			case tt.name == "balance successfully":
				s.EXPECT().GetUserBalance(gomock.Any(), tt.userID).Return(tt.want.current, tt.want.withdrawn, nil)
			case tt.name == "internal server error":
				var errServer = errors.New("internal server error")
				s.EXPECT().GetUserBalance(gomock.Any(), tt.userID).Return(float32(0), float32(0), errServer)
			default:
				s.EXPECT()
			}

			h := http.HandlerFunc(GetBalance(log, s))
			h(w, request)

			result := w.Result()
			body, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, string(body), tt.want.bodyResponse)
		})
	}
}
