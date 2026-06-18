package converter

import (
	"strings"

	"google.golang.org/protobuf/types/known/timestamppb"

	errs "github.com/anemptyemptiness/Go-Rocket-App/iam/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/model"
	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/service/input"
	commonv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/common/v1"
	userv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/user/v1"
)

func UserModelToProto(user model.User) *commonv1.User {
	pbUser := &commonv1.User{
		Uuid: user.UUID.String(),
		Info: &commonv1.UserInfo{
			Login: user.Login,
		},
		CreatedAt: timestamppb.New(user.CreatedAt),
	}

	if user.UpdatedAt != nil {
		pbUser.UpdatedAt = timestamppb.New(*user.UpdatedAt)
	}

	return pbUser
}

func RegisterRequestProtoToModel(req *userv1.RegisterRequest) (input.RegisterInput, error) {
	if req == nil || req.GetInfo() == nil || req.GetInfo().GetInfo() == nil {
		return input.RegisterInput{}, errs.ErrEmptyRequest
	}

	return input.RegisterInput{
		Info: model.UserRegistrationInfo{
			Info: model.UserInfo{
				Login: strings.TrimSpace(req.GetInfo().GetInfo().GetLogin()),
			},
			Password: strings.TrimSpace(req.GetInfo().GetPassword()),
		},
	}, nil
}

func GetUserRequestProtoToModel(req *userv1.GetUserRequest) (string, error) {
	if req == nil {
		return "", errs.ErrEmptyRequest
	}

	return req.GetUserUuid(), nil
}
