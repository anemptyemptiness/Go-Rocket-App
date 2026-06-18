package converter

import (
	"strings"

	"google.golang.org/protobuf/types/known/timestamppb"

	errs "github.com/anemptyemptiness/Go-Rocket-App/iam/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/model"
	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/service/input"
	authv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/auth/v1"
	commonv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/common/v1"
)

func SessionModelToProto(session model.Session) *commonv1.Session {
	pbSession := &commonv1.Session{
		CreatedAt: timestamppb.New(session.CreatedAt),
		ExpiresAt: timestamppb.New(session.ExpiresAt),
	}

	if session.UpdatedAt != nil {
		pbSession.UpdatedAt = timestamppb.New(*session.UpdatedAt)
	}

	return pbSession
}

func WhoamiRequestProtoToModel(req *authv1.WhoamiRequest) (string, error) {
	if req == nil {
		return "", errs.ErrEmptyRequest
	}

	return req.GetSessionUuid(), nil
}

func LoginRequestProtoToModel(req *authv1.LoginRequest) (input.LoginInput, error) {
	if req == nil {
		return input.LoginInput{}, errs.ErrEmptyRequest
	}

	return input.LoginInput{
		Login:    strings.TrimSpace(req.GetLogin()),
		Password: strings.TrimSpace(req.GetPassword()),
	}, nil
}

func LogoutRequestProtoToModel(req *authv1.LogoutRequest) (string, error) {
	if req == nil {
		return "", errs.ErrEmptyRequest
	}

	return req.GetSessionUuid(), nil
}
