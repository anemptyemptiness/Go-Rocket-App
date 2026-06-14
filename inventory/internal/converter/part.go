package converter

import (
	"strings"

	"google.golang.org/protobuf/types/known/timestamppb"

	errs "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/model/entity"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/model/valueobject"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/service/input"
	inventoryv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/inventory/v1"
)

func PartModelToProto(part entity.Part) *inventoryv1.Part {
	return &inventoryv1.Part{
		Uuid:          part.GetPartUUID(),
		Name:          part.GetName(),
		Description:   part.GetDescription(),
		Price:         part.GetPrice(),
		PartType:      valueobject.PartType.ToProto(part.GetPartType()),
		StockQuantity: part.GetStockQuantity(),
		CreatedAt:     timestamppb.New(part.GetCreatedAt()),
	}
}

func ListPartsRequestProtoToModel(req *inventoryv1.ListPartsRequest) (input.PartFilter, error) {
	if req == nil {
		return input.PartFilter{}, errs.ErrEmptyRequest
	}

	uuids := make([]string, 0, len(req.Uuids))
	for _, uuid := range req.Uuids {
		if uuid != "" {
			uuids = append(uuids, uuid)
		}
	}

	partType, err := valueobject.NewPartType(strings.TrimPrefix(inventoryv1.PartType_name[int32(req.GetPartType())], "PART_TYPE_"))
	if err != nil {
		return input.PartFilter{}, err
	}

	return input.PartFilter{
		UUIDs:    uuids,
		PartType: partType,
	}, nil
}

func PartsModelToProto(parts []entity.Part) []*inventoryv1.Part {
	protoParts := make([]*inventoryv1.Part, 0, len(parts))

	for _, part := range parts {
		protoParts = append(protoParts, PartModelToProto(part))
	}

	return protoParts
}

func GetPartRequestProtoToModel(req *inventoryv1.GetPartRequest) (string, error) {
	if req == nil {
		return "", errs.ErrEmptyRequest
	}

	return req.GetUuid(), nil
}

func ValidateCompatibilityRequestToModel(req *inventoryv1.ValidateCompatibilityRequest) (input.ValidateCompatibilityRequest, error) {
	if req == nil {
		return input.ValidateCompatibilityRequest{}, errs.ErrEmptyRequest
	}

	modelReq := input.ValidateCompatibilityRequest{
		HullUUID:   req.GetHullUuid(),
		EngineUUID: req.GetEngineUuid(),
	}

	if req.GetShieldUuid() != "" {
		modelReq.ShieldUUID = req.GetShieldUuid()
	}
	if req.GetWeaponUuid() != "" {
		modelReq.WeaponUUID = req.GetWeaponUuid()
	}

	return modelReq, nil
}

func ReservePartsRequestToModel(req *inventoryv1.ReservePartsRequest) (input.ReservePartsRequest, error) {
	if req == nil {
		return input.ReservePartsRequest{}, errs.ErrEmptyRequest
	}

	uuids := make([]string, 0, len(req.GetUuids()))
	for _, partUUID := range req.GetUuids() {
		if partUUID != "" {
			uuids = append(uuids, partUUID)
		}
	}

	return input.ReservePartsRequest{
		UUIDs: uuids,
	}, nil
}

func ReleasePartsRequestToModel(req *inventoryv1.ReleasePartsRequest) (input.ReleasePartsRequest, error) {
	if req == nil {
		return input.ReleasePartsRequest{}, errs.ErrEmptyRequest
	}

	uuids := make([]string, 0, len(req.GetUuids()))
	for _, partUUID := range req.GetUuids() {
		if partUUID != "" {
			uuids = append(uuids, partUUID)
		}
	}

	return input.ReleasePartsRequest{
		UUIDs: uuids,
	}, nil
}

func CommitPartsRequestToModel(req *inventoryv1.CommitPartsRequest) (input.CommitPartsRequest, error) {
	if req == nil {
		return input.CommitPartsRequest{}, errs.ErrEmptyRequest
	}

	uuids := make([]string, 0, len(req.GetUuids()))
	for _, partUUID := range req.GetUuids() {
		if partUUID != "" {
			uuids = append(uuids, partUUID)
		}
	}

	return input.CommitPartsRequest{
		UUIDs: uuids,
	}, nil
}
