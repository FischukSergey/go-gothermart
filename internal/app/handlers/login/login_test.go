package login

import (
	"bytes"
	"context"
	"github.com/FischukSergey/go-gothermart.git/internal/app/handlers/mock"
	"github.com/FischukSergey/go-gothermart.git/internal/app/middleware/auth"
	"github.com/FischukSergey/go-gothermart.git/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestLoginAuth(t *testing.T) {
	var log = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	type want struct {
		statusCode   int
		bodyResponse string
		user         models.User
		mockError    error
	}

	tests := []struct {
		name        string
		bodyRequest string
		want        want
	}{
		{
			name:        "login successfully",
			bodyRequest: `{"login":"family","password":"password"}`,
			want: want{
				statusCode:   200,
				bodyResponse: "",
				user: models.User{
					Email:    "family",
					Password: "password",
				},
			},
		},
		{
			name:        "password failed",
			bodyRequest: `{"login":"family","password":"password"}`,
			want: want{
				statusCode:   401,
				bodyResponse: "crypto/bcrypt: hashedPassword is not the hash of the given password" + "\n",
				user: models.User{
					Email:    "family",
					Password: "PASSWORD",
				},
			},
		},
		{
			name:        "password incorrect",
			bodyRequest: `{"login":"family","password":"pas"}`,
			want: want{
				statusCode:   401,
				bodyResponse: "Password: the length must be between 6 and 100." + "\n",
			},
		},
		{
			name:        "login incorrect",
			bodyRequest: `{"login":"fam","password":"password"}`,
			want: want{
				statusCode:   401,
				bodyResponse: "Email: the length must be between 6 and 100." + "\n",
			},
		},
		{
			name:        "internal server error",
			bodyRequest: `{"login":"family","password":"password"}`,
			want: want{
				statusCode:   500,
				bodyResponse: "can't create JWT, invalid user id or login" + "\n",
				user: models.User{
					Email:    "family",
					Password: "password",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			s := mock.NewMockLoginer(ctrl) //новый storage
			defer ctrl.Finish()

			requestTest := httptest.NewRequest(http.MethodPost, "/api/user/login",
				bytes.NewReader([]byte(tt.bodyRequest)))
			request := requestTest.WithContext(context.WithValue(requestTest.Context(),
				auth.CtxKeyUser, tt.want.user.ID))
			w := httptest.NewRecorder()
			request.Header.Set("Content-Type", "application/json")

			passHash, err := bcrypt.GenerateFromPassword([]byte(tt.want.user.Password), bcrypt.DefaultCost)
			if err != nil {
				log.Error("failed to generate password hash", "error", err)
				return
			}
			userResp := &models.User{
				EncryptedPassword: string(passHash),
				Email:             tt.want.user.Email,
				Role:              "",
				ID:                55,
			}

			switch {
			case tt.want.mockError != nil:
				s.EXPECT()
			case tt.name == "login successfully":
				s.EXPECT().Login(gomock.Any(), tt.want.user.Email).Return(userResp, nil)
			case tt.name == "password failed":
				s.EXPECT().Login(gomock.Any(), tt.want.user.Email).Return(userResp, tt.want.mockError)
			case tt.name == "password incorrect" || tt.name == "login incorrect":
				s.EXPECT()
			case tt.name == "internal server error":
				//var errServer = errors.New("internal server error")
				var u = &models.User{
					EncryptedPassword: string(passHash),
					Email:             "",
					ID:                0,
				}
				s.EXPECT().Login(gomock.Any(), tt.want.user.Email).Return(u, nil)
			default:
				s.EXPECT()
			}

			h := http.HandlerFunc(LoginAuth(log, s))
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
