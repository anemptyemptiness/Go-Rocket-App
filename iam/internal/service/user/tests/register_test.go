package tests

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	errs "github.com/anemptyemptiness/Go-Rocket-App/iam/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/model"
	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/service/input"
	usersvc "github.com/anemptyemptiness/Go-Rocket-App/iam/internal/service/user"
	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/service/user/mocks"
)

func Test_Register(t *testing.T) {
	t.Parallel()

	var (
		ctx           = context.Background()
		userUUID      = uuid.New()
		login         = gofakeit.Name()
		password      = "password!"
		shortPassword = "pass"
		unexpectedErr = assert.AnError
	)

	type args struct {
		req input.RegisterInput
	}

	type expected struct {
		err error
	}

	tests := []struct {
		name      string
		args      args
		expected  expected
		setupMock func(userRepo *mocks.UserRepository)
	}{
		{
			name: "успешная регистрация пользователя",
			args: args{
				req: input.RegisterInput{
					Info: model.UserRegistrationInfo{
						Info: model.UserInfo{
							Login: login,
						},
						Password: password,
					},
				},
			},
			expected: expected{
				err: nil,
			},
			setupMock: func(userRepo *mocks.UserRepository) {
				userRepo.EXPECT().
					GetByLogin(ctx, login).
					Return(model.User{}, errs.ErrUserNotFound)

				userRepo.EXPECT().
					Create(ctx, login, mock.MatchedBy(func(pwdHash string) bool {
						err := bcrypt.CompareHashAndPassword([]byte(pwdHash), []byte(password))
						if err != nil {
							return false
						}
						return true
					})).
					Return(userUUID, nil)
			},
		},
		{
			name: "ошибка: логин пустой",
			args: args{
				req: input.RegisterInput{
					Info: model.UserRegistrationInfo{
						Info: model.UserInfo{
							Login: "",
						},
						Password: password,
					},
				},
			},
			expected: expected{
				err: errs.ErrInvalidLogin,
			},
			setupMock: func(userRepo *mocks.UserRepository) {},
		},
		{
			name: "ошибка: пароль короткий",
			args: args{
				req: input.RegisterInput{
					Info: model.UserRegistrationInfo{
						Info: model.UserInfo{
							Login: login,
						},
						Password: shortPassword,
					},
				},
			},
			expected: expected{
				err: errs.ErrWeakPassword,
			},
			setupMock: func(userRepo *mocks.UserRepository) {},
		},
		{
			name: "ошибка: внутренняя ошибка user repository (GetByLogin)",
			args: args{
				req: input.RegisterInput{
					Info: model.UserRegistrationInfo{
						Info: model.UserInfo{
							Login: login,
						},
						Password: password,
					},
				},
			},
			expected: expected{
				err: unexpectedErr,
			},
			setupMock: func(userRepo *mocks.UserRepository) {
				userRepo.EXPECT().
					GetByLogin(ctx, login).
					Return(model.User{}, unexpectedErr)
			},
		},
		{
			name: "ошибка: пользователь уже существует",
			args: args{
				req: input.RegisterInput{
					Info: model.UserRegistrationInfo{
						Info: model.UserInfo{
							Login: login,
						},
						Password: password,
					},
				},
			},
			expected: expected{
				err: errs.ErrUserAlreadyExists,
			},
			setupMock: func(userRepo *mocks.UserRepository) {
				userRepo.EXPECT().
					GetByLogin(ctx, login).
					Return(model.User{}, errs.ErrUserAlreadyExists)
			},
		},
		{
			name: "ошибка: внутренняя ошибка user repository (Create)",
			args: args{
				req: input.RegisterInput{
					Info: model.UserRegistrationInfo{
						Info: model.UserInfo{
							Login: login,
						},
						Password: password,
					},
				},
			},
			expected: expected{
				err: unexpectedErr,
			},
			setupMock: func(userRepo *mocks.UserRepository) {
				userRepo.EXPECT().
					GetByLogin(ctx, login).
					Return(model.User{}, errs.ErrUserNotFound)

				userRepo.EXPECT().
					Create(ctx, login, mock.MatchedBy(func(pwdHash string) bool {
						err := bcrypt.CompareHashAndPassword([]byte(pwdHash), []byte(password))
						if err != nil {
							return false
						}
						return true
					})).
					Return(uuid.Nil, unexpectedErr)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			userRepo := mocks.NewUserRepository(t)

			tc.setupMock(userRepo)

			userSvc := usersvc.New(userRepo)

			userUuid, err := userSvc.Register(ctx, tc.args.req)
			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
				assert.Equal(t, uuid.Nil, userUuid)
			} else {
				require.NoError(t, err)
				assert.NotEqual(t, uuid.Nil, userUuid)
			}
		})
	}
}
