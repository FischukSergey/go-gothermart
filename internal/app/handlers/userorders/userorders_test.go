package userorders

import (
	"context"
	"errors"
	"github.com/FischukSergey/go-gothermart.git/internal/app/handlers/mock"
	"github.com/FischukSergey/go-gothermart.git/internal/app/middleware/auth"
	"github.com/FischukSergey/go-gothermart.git/internal/models"
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

func TestUserOrders(t *testing.T) {
	var log = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	type want struct {
		statusCode   int
		userID       int
		bodyResponse string
		mockError    error
	}

	tests := []struct {
		name   string
		userID int
		model  []models.GetUserOrders
		want   want
	}{
		{
			name:   "get orders successfully",
			userID: 55,
			model: []models.GetUserOrders{
				{
					Number:     "723893498573",
					Status:     "NEW",
					Accrual:    100.5,
					UploadedAt: "2024-08-17T17:01:40+03:00",
				},
			},
			want: want{
				statusCode: 200,
				bodyResponse: `[{"number":"723893498573","status":"NEW","accrual":100.5,` +
					`"uploaded_at":"2024-08-17T17:01:40+03:00"}]` + "\n",
				userID: 1,
			},
		},
		{
			name:   "not found orders",
			userID: 55,
			model:  []models.GetUserOrders{},
			want: want{
				statusCode:   204,
				bodyResponse: `"not found data for user"` + "\n",
				userID:       1,
			},
		},
		{
			name:   "internal server error",
			userID: 55,
			model:  []models.GetUserOrders{},
			want: want{
				statusCode:   500,
				bodyResponse: `{"error":"internal server error"}` + "\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			s := mock.NewMockUserOrdersGetter(ctrl) //новый storage
			defer ctrl.Finish()

			requestTest := httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)
			request := requestTest.WithContext(context.WithValue(requestTest.Context(),
				auth.CtxKeyUser, tt.userID))
			w := httptest.NewRecorder()

			switch {
			case tt.want.mockError != nil:
				s.EXPECT()
			case tt.name == "get orders successfully":
				s.EXPECT().GetUserOrders(gomock.Any(), tt.userID).Return(tt.model, nil)
			case tt.name == "not found orders":
				s.EXPECT().GetUserOrders(gomock.Any(), gomock.Any()).Return(tt.model, nil)
			case tt.name == "internal server error":
				var errServer = errors.New("internal server error")
				s.EXPECT().GetUserOrders(gomock.Any(), gomock.Any()).Return(tt.model, errServer)
			default:
				s.EXPECT()
			}

			h := http.HandlerFunc(UserOrders(log, s))
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
