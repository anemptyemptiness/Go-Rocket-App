package converter

import (
	"github.com/google/uuid"

	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/model"
	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/repository/record"
)

func UserModelToRecord(user model.User) record.User {
	return record.User{
		UUID:         user.UUID.String(),
		Login:        user.Login,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}

func UserRecordToModel(user record.User) model.User {
	return model.User{
		UUID:         uuid.MustParse(user.UUID),
		Login:        user.Login,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}
