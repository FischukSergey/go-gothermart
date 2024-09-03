package withdraw

import (
	"bytes"
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

func TestOrderWithdraw(t *testing.T) {
	var log = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	type want struct {
		statusCode   int
		bodyResponse string
		mockError    error
	}

	tests := []struct {
		name        string
		bodyRequest string
		model       models.Order
		want        want
	}{
		{
			name:        "withdraw successfully",
			bodyRequest: `{"order":"723893498573","sum":1000}`,
			model: models.Order{
				UserID:    55,
				OrderID:   "723893498573",
				Withdraw:  1000,
				Status:    "",
				CreatedAt: time.Now().Format(time.RFC3339),
			},
			want: want{
				statusCode:   200,
				bodyResponse: "",
			},
		},
		{
			name:        "invalid order",
			bodyRequest: `{"order":"723893498575","sum":1000}`,
			model: models.Order{
				UserID:    55,
				OrderID:   "723893498573",
				Withdraw:  1000,
				Status:    "",
				CreatedAt: time.Now().Format(time.RFC3339),
			},
			want: want{
				statusCode:   422,
				bodyResponse: "Number order is invalid" + "\n",
			},
		},
		{
			name:        "insufficient amount",
			bodyRequest: `{"order":"723893498573","sum":1000}`,
			model: models.Order{
				UserID:    55,
				OrderID:   "723893498573",
				Withdraw:  1000,
				CreatedAt: time.Now().Format(time.RFC3339),
			},
			want: want{
				statusCode:   402,
				bodyResponse: "insufficient funds" + "\n",
			},
		},
		{
			name:        "internal server error",
			bodyRequest: `{"order":"723893498573","sum":1000}`,
			model: models.Order{
				UserID:    55,
				OrderID:   "723893498573",
				Withdraw:  1000,
				CreatedAt: time.Now().Format(time.RFC3339),
			},
			want: want{
				statusCode:   500,
				bodyResponse: "internal server error" + "\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			s := mock.NewMockOrderBalanceWithdraw(ctrl) //новый storage
			defer ctrl.Finish()

			requestTest := httptest.NewRequest(http.MethodPost, "/api/balance/withdraw",
				bytes.NewReader([]byte(tt.bodyRequest)))
			request := requestTest.WithContext(context.WithValue(requestTest.Context(),
				auth.CtxKeyUser, tt.model.UserID))
			w := httptest.NewRecorder()
			request.Header.Set("Content-Type", "application/json")

			switch {
			case tt.want.mockError != nil:
				s.EXPECT()
			case tt.name == "withdraw successfully":
				s.EXPECT().CreateOrderWithdraw(gomock.Any(), tt.model).Return(nil)
			case tt.name == "insufficient amount":
				s.EXPECT().CreateOrderWithdraw(gomock.Any(), tt.model).Return(models.ErrInsufficientFunds)
			case tt.name == "internal server error":
				var errServer = errors.New("internal server error")
				s.EXPECT().CreateOrderWithdraw(gomock.Any(), gomock.Any()).Return(errServer)
			default:
				s.EXPECT()
			}

			h := http.HandlerFunc(OrderWithdraw(log, s))
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
