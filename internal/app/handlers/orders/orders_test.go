package orders

import (
	"bytes"
	"context"
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

func TestOrderSave(t *testing.T) {
	var log = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	type want struct {
		statusCode   int
		bodyResponse string
		errError     error
		mockError    error
	}

	tests := []struct {
		name        string
		bodyRequest string
		userID      int
		want        want
	}{
		{
			name:        "save order successfully",
			bodyRequest: `723893498573`,
			userID:      1,
			want: want{
				statusCode:   202,
				bodyResponse: ``,
			},
		},
		{
			name:        "empty request",
			bodyRequest: ``,
			userID:      1,
			want: want{
				statusCode:   400,
				bodyResponse: "Request is empty" + "\n",
			},
		},
		{
			name:        "invalid order",
			bodyRequest: "khjkj",
			userID:      1,
			want: want{
				statusCode:   422,
				bodyResponse: "Request is invalid" + "\n",
			},
		},
		{
			name:        "order already exists",
			bodyRequest: "723893498573",
			userID:      1,
			want: want{
				statusCode:   200,
				bodyResponse: "order already loaded same user" + "\n",
				errError:     models.ErrOrderUploadedSameUser,
			},
		},
		{
			name:        "order already exists another user",
			bodyRequest: "723893498573",
			userID:      1,
			want: want{
				statusCode:   409,
				bodyResponse: "order loaded another user" + "\n",
				errError:     models.ErrOrderUploadedAnotherUser,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			s := mock.NewMockOrderSaver(ctrl) //новый storage
			defer ctrl.Finish()

			u := models.Order{ //модель ордера для передачи в БД(мок)
				UserID:    1,
				OrderID:   tt.bodyRequest,
				Status:    "NEW",
				CreatedAt: time.Now().Format(time.RFC3339),
			}

			requestTest := httptest.NewRequest(http.MethodPost,
				"/api/user/orders",
				bytes.NewReader([]byte(tt.bodyRequest)))
			request := requestTest.WithContext(context.WithValue(requestTest.Context(),
				auth.CtxKeyUser, tt.userID))
			w := httptest.NewRecorder()
			request.Header.Set("Content-Type", "text/plain")

			switch {
			case tt.want.mockError != nil:
				s.EXPECT()
			case tt.name == "save order successfully":
				s.EXPECT().CreateOrder(gomock.Any(), u).Return(nil)
			case tt.name == "order already exists":
				s.EXPECT().CreateOrder(gomock.Any(), u).Return(tt.want.errError)
			case tt.name == "order already exists another user":
				s.EXPECT().CreateOrder(gomock.Any(), u).Return(tt.want.errError)
			default:
				s.EXPECT()
			}

			h := http.HandlerFunc(OrderSave(log, s))
			h(w, request)

			result := w.Result()
			body, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.bodyResponse, string(body))
		})
	}
}
