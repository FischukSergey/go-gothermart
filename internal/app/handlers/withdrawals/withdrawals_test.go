package withdrawals

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
	"time"
)

func TestOrderWithdrawAll(t *testing.T) {
	var log = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	type want struct {
		statusCode   int
		bodyResponse string
		resp         []models.GetAllWithdraw
		mockError    error
	}

	tests := []struct {
		name   string
		userID int
		want   want
	}{
		{
			name:   "withdrawals successfully",
			userID: 55,
			want: want{
				statusCode: 200,
				bodyResponse: `[{` +
					`"order":"723893498573","sum":100,"processed_at":"` + time.Now().Format(time.RFC3339) + `"},{` +
					`"order":"465893498571","sum":150,"processed_at":"` + time.Now().Format(time.RFC3339) + `"}]` + "\n",
				resp: []models.GetAllWithdraw{
					{
						Order:       "723893498573",
						Sum:         100,
						ProcessedAt: time.Now().Format(time.RFC3339),
					},
					{
						Order:       "465893498571",
						Sum:         150,
						ProcessedAt: time.Now().Format(time.RFC3339),
					},
				},
			},
		},
		{
			name:   "not found withdraw orders for user",
			userID: 55,
			want: want{
				statusCode:   204,
				resp:         []models.GetAllWithdraw{},
				bodyResponse: `"not found withdraw orders for user"` + "\n",
			},
		},
		{
			name:   "internal server error",
			userID: 55,
			want: want{
				statusCode:   500,
				bodyResponse: "internal server error" + "\n",
				resp:         nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			s := mock.NewMockGetUserWithdrawAll(ctrl) //новый storage
			defer ctrl.Finish()

			requestTest := httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", nil)
			request := requestTest.WithContext(context.WithValue(requestTest.Context(),
				auth.CtxKeyUser, tt.userID))
			w := httptest.NewRecorder()
			request.Header.Set("Content-Type", "text/plain")

			switch {
			case tt.want.mockError != nil:
				s.EXPECT()
			case tt.name == "withdrawals successfully":
				s.EXPECT().GetAllWithdraw(gomock.Any(), tt.userID).Return(tt.want.resp, nil)
			case tt.name == "not found withdraw orders for user":
				s.EXPECT().GetAllWithdraw(gomock.Any(), tt.userID).Return(tt.want.resp, tt.want.mockError)
			case tt.name == "internal server error":
				var errServer = errors.New("internal server error")
				s.EXPECT().GetAllWithdraw(gomock.Any(), tt.userID).Return(tt.want.resp, errServer)
			default:
				s.EXPECT()
			}

			h := http.HandlerFunc(OrderWithdrawAll(log, s))
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
