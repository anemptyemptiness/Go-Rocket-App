package tests

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	errs "github.com/anemptyemptiness/Go-Rocket-App/iam/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/model"
	authsvc "github.com/anemptyemptiness/Go-Rocket-App/iam/internal/service/auth"
	authmocks "github.com/anemptyemptiness/Go-Rocket-App/iam/internal/service/auth/mocks"
	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/service/input"
	usermocks "github.com/anemptyemptiness/Go-Rocket-App/iam/internal/service/user/mocks"
)

func Test_Login(t *testing.T) {
	t.Parallel()

	var (
		ctx             = context.Background()
		login           = gofakeit.Name()
		userUUID        = uuid.New()
		sessionUUID     = uuid.New()
		now             = time.Now()
		passwordValid   = "passwordXXX"
		passwordInvalid = "passwor"
		unexpectedErr   = assert.AnError
	)

	passwordValidHash, err := bcrypt.GenerateFromPassword([]byte(passwordValid), bcrypt.MinCost)
	require.NoError(t, err)

	user := model.User{
		UUID:         userUUID,
		Login:        login,
		PasswordHash: string(passwordValidHash),
		CreatedAt:    now,
		UpdatedAt:    &now,
	}

	type args struct {
		req input.LoginInput
	}

	type expected struct {
		err error
	}

	tests := []struct {
		name      string
		args      args
		expected  expected
		setupMock func(sessionRepo *authmocks.SessionRepository, userRepo *usermocks.UserRepository)
	}{
		{
			name: "успешный логин пользователя",
			args: args{
				req: input.LoginInput{
					Login:    login,
					Password: passwordValid,
				},
			},
			expected: expected{
				err: nil,
			},
			setupMock: func(sessionRepo *authmocks.SessionRepository, userRepo *usermocks.UserRepository) {
				userRepo.EXPECT().
					GetByLogin(ctx, login).
					Return(user, nil)

				sessionRepo.EXPECT().
					Create(ctx, user).
					Return(sessionUUID, nil)
			},
		},
		{
			name: "ошибка: пустой логин",
			args: args{
				req: input.LoginInput{
					Login:    "",
					Password: passwordValid,
				},
			},
			expected: expected{
				err: errs.ErrEmptyCredential,
			},
			setupMock: func(sessionRepo *authmocks.SessionRepository, userRepo *usermocks.UserRepository) {},
		},
		{
			name: "ошибка: пустой пароль",
			args: args{
				req: input.LoginInput{
					Login:    login,
					Password: "",
				},
			},
			expected: expected{
				err: errs.ErrEmptyCredential,
			},
			setupMock: func(sessionRepo *authmocks.SessionRepository, userRepo *usermocks.UserRepository) {},
		},
		{
			name: "ошибка: пользователь не найден",
			args: args{
				req: input.LoginInput{
					Login:    login,
					Password: passwordValid,
				},
			},
			expected: expected{
				err: errs.ErrUserNotFound,
			},
			setupMock: func(sessionRepo *authmocks.SessionRepository, userRepo *usermocks.UserRepository) {
				userRepo.EXPECT().
					GetByLogin(ctx, login).
					Return(model.User{}, errs.ErrUserNotFound)
			},
		},
		{
			name: "ошибка: внутренняя ошибка user repository",
			args: args{
				req: input.LoginInput{
					Login:    login,
					Password: passwordValid,
				},
			},
			expected: expected{
				err: unexpectedErr,
			},
			setupMock: func(sessionRepo *authmocks.SessionRepository, userRepo *usermocks.UserRepository) {
				userRepo.EXPECT().
					GetByLogin(ctx, login).
					Return(model.User{}, unexpectedErr)
			},
		},
		{
			name: "ошибка: пароли не совпадают",
			args: args{
				req: input.LoginInput{
					Login:    login,
					Password: passwordInvalid,
				},
			},
			expected: expected{
				err: errs.ErrInvalidCredentials,
			},
			setupMock: func(sessionRepo *authmocks.SessionRepository, userRepo *usermocks.UserRepository) {
				userRepo.EXPECT().
					GetByLogin(ctx, login).
					Return(user, nil)
			},
		},
		{
			name: "ошибка: внутренняя ошибка session repository",
			args: args{
				req: input.LoginInput{
					Login:    login,
					Password: passwordValid,
				},
			},
			expected: expected{
				err: unexpectedErr,
			},
			setupMock: func(sessionRepo *authmocks.SessionRepository, userRepo *usermocks.UserRepository) {
				userRepo.EXPECT().
					GetByLogin(ctx, login).
					Return(user, nil)

				sessionRepo.EXPECT().
					Create(ctx, user).
					Return(uuid.Nil, unexpectedErr)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			authRepo := authmocks.NewSessionRepository(t)
			userRepo := usermocks.NewUserRepository(t)
			tc.setupMock(authRepo, userRepo)

			authSvc := authsvc.New(authRepo, userRepo)

			sessionUuid, err := authSvc.Login(ctx, tc.args.req)
			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
			} else {
				require.NoError(t, err)
				assert.NotEqual(t, uuid.Nil, sessionUuid)
			}
		})
	}
}
