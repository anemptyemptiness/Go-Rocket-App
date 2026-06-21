package converter

import (
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	commonv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/common/v1"
)

func SessionProtoToModel(session *commonv1.Session) model.Session {
	if session == nil {
		return model.Session{}
	}

	return model.Session{
		UUID:      session.GetUuid(),
		CreatedAt: session.GetCreatedAt().AsTime(),
		UpdatedAt: session.GetUpdatedAt().AsTime(),
		ExpiresAt: session.GetExpiresAt().AsTime(),
	}
}

func UserProtoToModel(user *commonv1.User) model.User {
	if user == nil {
		return model.User{}
	}

	var login string
	if user.GetInfo() != nil {
		login = user.GetInfo().GetLogin()
	}

	return model.User{
		UUID: user.GetUuid(),
		Info: model.UserInfo{
			Login: login,
		},
		CreatedAt: user.GetCreatedAt().AsTime(),
		UpdatedAt: user.GetUpdatedAt().AsTime(),
	}
}
