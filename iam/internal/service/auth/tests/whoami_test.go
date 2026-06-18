package tests

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	errs "github.com/anemptyemptiness/Go-Rocket-App/iam/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/model"
	authsvc "github.com/anemptyemptiness/Go-Rocket-App/iam/internal/service/auth"
	authmocks "github.com/anemptyemptiness/Go-Rocket-App/iam/internal/service/auth/mocks"
	usermocks "github.com/anemptyemptiness/Go-Rocket-App/iam/internal/service/user/mocks"
)

func Test_Whoami(t *testing.T) {
	t.Parallel()

	var (
		ctx           = context.Background()
		sessionUUID   = gofakeit.UUID()
		userUUID      = gofakeit.UUID()
		login         = gofakeit.Name()
		now           = time.Now()
		nowWithTTL    = now.Add(24 * time.Hour)
		unexpectedErr = assert.AnError
	)

	session := model.Session{
		UUID:      uuid.MustParse(sessionUUID),
		UserUUID:  uuid.MustParse(userUUID),
		Login:     login,
		CreatedAt: now,
		UpdatedAt: &now,
		ExpiresAt: nowWithTTL,
	}

	user := model.User{
		UUID:      uuid.MustParse(userUUID),
		Login:     login,
		CreatedAt: now,
		UpdatedAt: &now,
	}

	type args struct {
		sessionUuid string
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
			name: "успешная идентификация сессии и пользователя",
			args: args{
				sessionUuid: sessionUUID,
			},
			expected: expected{
				err: nil,
			},
			setupMock: func(sessionRepo *authmocks.SessionRepository, userRepo *usermocks.UserRepository) {
				sessionRepo.EXPECT().
					GetByUUID(ctx, mock.MatchedBy(func(sessionUuid uuid.UUID) bool {
						return sessionUuid != uuid.Nil &&
							sessionUuid == uuid.MustParse(sessionUUID)
					})).
					Return(session, nil)

				userRepo.EXPECT().
					GetByUUID(ctx, session.UserUUID).
					Return(user, nil)
			},
		},
		{
			name: "ошибка: uuid сессии невалидный",
			args: args{
				sessionUuid: "fdsfndofnsdfk",
			},
			expected: expected{
				err: errs.ErrInvalidUUID,
			},
			setupMock: func(sessionRepo *authmocks.SessionRepository, userRepo *usermocks.UserRepository) {},
		},
		{
			name: "ошибка: сессия не найдена",
			args: args{
				sessionUuid: sessionUUID,
			},
			expected: expected{
				err: errs.ErrSessionNotFound,
			},
			setupMock: func(sessionRepo *authmocks.SessionRepository, userRepo *usermocks.UserRepository) {
				sessionRepo.EXPECT().
					GetByUUID(ctx, mock.MatchedBy(func(sessionUuid uuid.UUID) bool {
						return sessionUuid != uuid.Nil &&
							sessionUuid == uuid.MustParse(sessionUUID)
					})).
					Return(model.Session{}, errs.ErrSessionNotFound)
			},
		},
		{
			name: "ошибка: внутренняя ошибка session repository",
			args: args{
				sessionUuid: sessionUUID,
			},
			expected: expected{
				err: unexpectedErr,
			},
			setupMock: func(sessionRepo *authmocks.SessionRepository, userRepo *usermocks.UserRepository) {
				sessionRepo.EXPECT().
					GetByUUID(ctx, mock.MatchedBy(func(sessionUuid uuid.UUID) bool {
						return sessionUuid != uuid.Nil &&
							sessionUuid == uuid.MustParse(sessionUUID)
					})).
					Return(model.Session{}, unexpectedErr)
			},
		},
		{
			name: "ошибка: пользователь не найден",
			args: args{
				sessionUuid: sessionUUID,
			},
			expected: expected{
				err: errs.ErrUserNotFound,
			},
			setupMock: func(sessionRepo *authmocks.SessionRepository, userRepo *usermocks.UserRepository) {
				sessionRepo.EXPECT().
					GetByUUID(ctx, mock.MatchedBy(func(sessionUuid uuid.UUID) bool {
						return sessionUuid != uuid.Nil &&
							sessionUuid == uuid.MustParse(sessionUUID)
					})).
					Return(session, nil)

				userRepo.EXPECT().
					GetByUUID(ctx, session.UserUUID).
					Return(model.User{}, errs.ErrUserNotFound)
			},
		},
		{
			name: "ошибка: внутренняя ошибка user repository",
			args: args{
				sessionUuid: sessionUUID,
			},
			expected: expected{
				err: unexpectedErr,
			},
			setupMock: func(sessionRepo *authmocks.SessionRepository, userRepo *usermocks.UserRepository) {
				sessionRepo.EXPECT().
					GetByUUID(ctx, mock.MatchedBy(func(sessionUuid uuid.UUID) bool {
						return sessionUuid != uuid.Nil &&
							sessionUuid == uuid.MustParse(sessionUUID)
					})).
					Return(session, nil)

				userRepo.EXPECT().
					GetByUUID(ctx, session.UserUUID).
					Return(model.User{}, unexpectedErr)
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

			respSession, respUser, err := authSvc.Whoami(ctx, tc.args.sessionUuid)
			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, session.UUID, respSession.UUID)
				assert.Equal(t, user.UUID, respUser.UUID)
			}
		})
	}
}
