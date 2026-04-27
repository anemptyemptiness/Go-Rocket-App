package app

import (
	"context"
	"sort"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	inventoryapi "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/api/inventory/v1"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/interceptor"
	partrepo "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/repository/part"
	partsvc "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/service/part"
	inventoryv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/inventory/v1"
)

// Part представляет деталь космического корабля.
type Part struct {
	UUID          string
	Name          string
	Description   string
	Price         int64 // в копейках
	PartType      inventoryv1.PartType
	StockQuantity int64
	CreatedAt     *timestamppb.Timestamp
}

// InventoryServer реализует gRPC сервис.
type InventoryServer struct {
	inventoryv1.UnimplementedInventoryServiceServer
	parts map[uuid.UUID]Part
}

// NewInventoryServer создаёт сервер с предзагруженными seed-данными.
func NewInventoryServer() *InventoryServer {
	now := timestamppb.Now()

	return &InventoryServer{
		parts: map[uuid.UUID]Part{
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440001",
				Name:          "Алюминиевый корпус",
				Description:   "Лёгкий корпус для небольших кораблей",
				Price:         500000, // 5000₽
				PartType:      inventoryv1.PartType_PART_TYPE_HULL,
				StockQuantity: 10,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440002",
				Name:          "Титановый корпус",
				Description:   "Прочный корпус для средних кораблей",
				Price:         1500000, // 15000₽
				PartType:      inventoryv1.PartType_PART_TYPE_HULL,
				StockQuantity: 5,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440003"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440003",
				Name:          "Ионный двигатель C",
				Description:   "Базовый ионный двигатель класса C",
				Price:         300000, // 3000₽
				PartType:      inventoryv1.PartType_PART_TYPE_ENGINE,
				StockQuantity: 8,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440004"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440004",
				Name:          "Ионный двигатель B",
				Description:   "Улучшенный ионный двигатель класса B",
				Price:         800000, // 8000₽
				PartType:      inventoryv1.PartType_PART_TYPE_ENGINE,
				StockQuantity: 3,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440005"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440005",
				Name:          "Энергетический щит",
				Description:   "Стандартный энергетический щит",
				Price:         400000, // 4000₽
				PartType:      inventoryv1.PartType_PART_TYPE_SHIELD,
				StockQuantity: 6,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440006"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440006",
				Name:          "Лазерная пушка",
				Description:   "Точная лазерная пушка",
				Price:         250000, // 2500₽
				PartType:      inventoryv1.PartType_PART_TYPE_WEAPON,
				StockQuantity: 7,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440007"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440007",
				Name:          "Очень полезная вещь",
				Description:   "Описание очень полезной вещи",
				Price:         2000000, // 20000₽
				PartType:      inventoryv1.PartType_PART_TYPE_HULL,
				StockQuantity: 0,
				CreatedAt:     now,
			},
		},
	}
}

func RegisterServices(grpcServer *grpc.Server) {
	repo := partrepo.New()
	svc := partsvc.New(repo)
	api := inventoryapi.New(svc)
	inventoryv1.RegisterInventoryServiceServer(grpcServer, api)
}

func Interceptors() []grpc.ServerOption {
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptor.ErrorInterceptor(),
		),
	}

	return opts
}

// GetPart возвращает деталь по UUID.
func (s *InventoryServer) GetPart(
	ctx context.Context,
	req *inventoryv1.GetPartRequest,
) (*inventoryv1.GetPartResponse, error) {
	if req.GetUuid() == "" {
		return nil, status.Error(codes.InvalidArgument, "идентификатор детали не может быть пустым")
	}

	partUUID, err := uuid.Parse(req.GetUuid())
	if err != nil || partUUID == uuid.Nil {
		return nil, status.Errorf(codes.InvalidArgument, "идентификатор детали невалидный UUID: %v", err)
	}

	part, ok := s.parts[partUUID]
	if !ok {
		return nil, status.Error(codes.NotFound, "деталь не найдена")
	}

	pPart := &inventoryv1.Part{
		Uuid:          partUUID.String(),
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		PartType:      part.PartType,
		StockQuantity: part.StockQuantity,
		CreatedAt:     part.CreatedAt,
	}

	return &inventoryv1.GetPartResponse{
		Part: pPart,
	}, nil
}

// ListParts возвращает список деталей с опциональной фильтрацией по типу.
func (s *InventoryServer) ListParts(
	ctx context.Context,
	req *inventoryv1.ListPartsRequest,
) (*inventoryv1.ListPartsResponse, error) {
	parts := make([]*inventoryv1.Part, 0, len(req.GetUuids()))
	if len(req.GetUuids()) > 0 {
		for _, stringUUID := range req.GetUuids() {
			parsedUUID, err := uuid.Parse(stringUUID)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "невалидный UUID детали: %v", err)
			}

			part, exists := s.parts[parsedUUID]
			if !exists {
				return nil, status.Error(codes.NotFound, "деталь не найдена")
			}

			parts = append(parts, &inventoryv1.Part{
				Uuid:          parsedUUID.String(),
				Name:          part.Name,
				Description:   part.Description,
				Price:         part.Price,
				PartType:      part.PartType,
				StockQuantity: part.StockQuantity,
				CreatedAt:     part.CreatedAt,
			})
		}
	} else {
		for _, part := range s.parts {
			if req.GetPartType() == inventoryv1.PartType_PART_TYPE_UNSPECIFIED || req.GetPartType() == part.PartType {
				parts = append(parts, &inventoryv1.Part{
					Uuid:          part.UUID,
					Name:          part.Name,
					Description:   part.Description,
					Price:         part.Price,
					PartType:      part.PartType,
					StockQuantity: part.StockQuantity,
					CreatedAt:     part.CreatedAt,
				})
			}
		}

		sort.Slice(parts, func(i, j int) bool {
			return parts[i].Name < parts[j].Name
		})
	}

	return &inventoryv1.ListPartsResponse{
		Parts: parts,
	}, nil
}
