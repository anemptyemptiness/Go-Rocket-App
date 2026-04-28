package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	inventoryapi "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/api/inventory/v1"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/api/inventory/v1/mocks"
	errs "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/model"
	inventoryv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/inventory/v1"
)

func TestListParts(t *testing.T) {
	t.Parallel()

	type args struct {
		req *inventoryv1.ListPartsRequest
	}

	type expected struct {
		resp    *inventoryv1.ListPartsResponse
		wantErr error
	}

	var (
		ctx = context.Background()

		hullUUID            = gofakeit.UUID()
		hullUUIDInvalid     = "fapskfpakf"
		engineUUID          = gofakeit.UUID()
		shieldUUID          = gofakeit.UUID()
		weaponUUID          = gofakeit.UUID()
		partTypeUnspecified = inventoryv1.PartType_PART_TYPE_UNSPECIFIED
		partTypeHull        = inventoryv1.PartType_PART_TYPE_HULL
		partTypeEngine      = inventoryv1.PartType_PART_TYPE_ENGINE
		partTypeShield      = inventoryv1.PartType_PART_TYPE_SHIELD
		partTypeWeapon      = inventoryv1.PartType_PART_TYPE_WEAPON

		internalError = errors.New("неожиданность")

		now = time.Now()
	)

	tests := []struct {
		name      string
		args      args
		expected  expected
		setupMock func(svc *mocks.InventoryService)
	}{
		{
			name: "успешное получение деталей",
			args: args{
				req: &inventoryv1.ListPartsRequest{
					PartType: partTypeUnspecified,
					Uuids:    []string{shieldUUID, engineUUID, weaponUUID, hullUUID},
				},
			},
			expected: expected{
				resp: &inventoryv1.ListPartsResponse{
					Parts: []*inventoryv1.Part{
						{
							Uuid:          shieldUUID,
							Name:          "shield",
							Description:   "shield-desc",
							Price:         10000,
							PartType:      partTypeShield,
							StockQuantity: 10,
							CreatedAt:     timestamppb.New(now),
						},
						{
							Uuid:          engineUUID,
							Name:          "engine",
							Description:   "engine-desc",
							Price:         20000,
							PartType:      partTypeEngine,
							StockQuantity: 20,
							CreatedAt:     timestamppb.New(now),
						},
						{
							Uuid:          weaponUUID,
							Name:          "weapon",
							Description:   "weapon-desc",
							Price:         30000,
							PartType:      partTypeWeapon,
							StockQuantity: 30,
							CreatedAt:     timestamppb.New(now),
						},
						{
							Uuid:          hullUUID,
							Name:          "hull",
							Description:   "hull-desc",
							Price:         40000,
							PartType:      partTypeHull,
							StockQuantity: 40,
							CreatedAt:     timestamppb.New(now),
						},
					},
				},
				wantErr: nil,
			},
			setupMock: func(svc *mocks.InventoryService) {
				svc.EXPECT().
					ListParts(ctx, []string{shieldUUID, engineUUID, weaponUUID, hullUUID}, model.PartTypeUnspecified).
					Return([]model.Part{
						{
							UUID:          shieldUUID,
							Name:          "shield",
							Description:   "shield-desc",
							Price:         10000,
							PartType:      model.PartTypeShield,
							StockQuantity: 10,
							CreatedAt:     now,
						},
						{
							UUID:          engineUUID,
							Name:          "engine",
							Description:   "engine-desc",
							Price:         20000,
							PartType:      model.PartTypeEngine,
							StockQuantity: 20,
							CreatedAt:     now,
						},
						{
							UUID:          weaponUUID,
							Name:          "weapon",
							Description:   "weapon-desc",
							Price:         30000,
							PartType:      model.PartTypeWeapon,
							StockQuantity: 30,
							CreatedAt:     now,
						},
						{
							UUID:          hullUUID,
							Name:          "hull",
							Description:   "hull-desc",
							Price:         40000,
							PartType:      model.PartTypeHull,
							StockQuantity: 40,
							CreatedAt:     now,
						},
					}, nil)
			},
		},
		{
			name: "ошибка: пустой реквест",
			args: args{
				req: nil,
			},
			expected: expected{
				resp:    nil,
				wantErr: errs.ErrEmptyRequest,
			},
			setupMock: func(svc *mocks.InventoryService) {},
		},
		{
			name: "ошибка: невалидный UUID детали",
			args: args{
				req: &inventoryv1.ListPartsRequest{
					Uuids:    []string{hullUUIDInvalid},
					PartType: partTypeUnspecified,
				},
			},
			expected: expected{
				resp:    nil,
				wantErr: errs.ErrIncorrectPartUUID,
			},
			setupMock: func(svc *mocks.InventoryService) {
				svc.EXPECT().
					ListParts(ctx, []string{hullUUIDInvalid}, model.PartTypeUnspecified).
					Return(nil, errs.ErrIncorrectPartUUID)
			},
		},
		{
			name: "ошибка: деталь не найдена",
			args: args{
				req: &inventoryv1.ListPartsRequest{
					Uuids:    []string{engineUUID},
					PartType: partTypeEngine,
				},
			},
			expected: expected{
				resp:    nil,
				wantErr: errs.ErrPartNotFound,
			},
			setupMock: func(svc *mocks.InventoryService) {
				svc.EXPECT().
					ListParts(ctx, []string{engineUUID}, model.PartTypeEngine).
					Return(nil, errs.ErrPartNotFound)
			},
		},
		{
			name: "ошибка: внутренняя ошибка",
			args: args{
				req: &inventoryv1.ListPartsRequest{
					Uuids:    []string{weaponUUID},
					PartType: partTypeWeapon,
				},
			},
			expected: expected{
				resp:    nil,
				wantErr: internalError,
			},
			setupMock: func(svc *mocks.InventoryService) {
				svc.EXPECT().
					ListParts(ctx, []string{weaponUUID}, model.PartTypeWeapon).
					Return(nil, internalError)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := mocks.NewInventoryService(t)
			tc.setupMock(svc)

			api := inventoryapi.New(svc)

			resp, err := api.ListParts(ctx, tc.args.req)
			if tc.expected.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.wantErr)
				assert.Empty(t, resp)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, resp)
				assert.IsType(t, &inventoryv1.ListPartsResponse{}, resp)
				assert.Len(t, resp.GetParts(), len(tc.expected.resp.GetParts()))

				for idx, part := range resp.GetParts() {
					require.NotNil(t, part)
					require.Less(t, idx, len(tc.expected.resp.GetParts()))
					assert.Equal(t, tc.expected.resp.GetParts()[idx].GetUuid(), part.GetUuid())
				}
			}
		})
	}
}
