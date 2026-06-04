package tests

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errs "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/model/entity"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/model/valueobject"
	compatibilitychecker "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/service/domain/compatibility_checker"
	pkgerr "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/errors"
)

func TestCheck(t *testing.T) {
	t.Parallel()

	type args struct {
		parts []entity.Part
	}

	type expected struct {
		wantErr error
	}

	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "compatible hull and engine class C",
			args: args{
				parts: []entity.Part{
					newHullPart(t, 31),
					newEnginePart(t, valueobject.EngineClassC, 30),
				},
			},
			expected: expected{
				wantErr: nil,
			},
		},
		{
			name: "compatible hull and engine class B",
			args: args{
				parts: []entity.Part{
					newHullPart(t, 70),
					newEnginePart(t, valueobject.EngineClassB, 70),
				},
			},
			expected: expected{
				wantErr: nil,
			},
		},
		{
			name: "incompatible hull and engine class B",
			args: args{
				parts: []entity.Part{
					newHullPart(t, 69),
					newEnginePart(t, valueobject.EngineClassB, 70),
				},
			},
			expected: expected{
				wantErr: pkgerr.FailedPrecondition(errs.ErrIncompatibleParts),
			},
		},
		{
			name: "compatible hull and engine class A",
			args: args{
				parts: []entity.Part{
					newHullPart(t, 100),
					newEnginePart(t, valueobject.EngineClassA, 100),
				},
			},
			expected: expected{
				wantErr: nil,
			},
		},
		{
			name: "incompatible hull and engine class A",
			args: args{
				parts: []entity.Part{
					newHullPart(t, 99),
					newEnginePart(t, valueobject.EngineClassA, 100),
				},
			},
			expected: expected{
				wantErr: pkgerr.FailedPrecondition(errs.ErrIncompatibleParts),
			},
		},
		{
			name: "compatible energy shield and laser weapon",
			args: args{
				parts: []entity.Part{
					newHullPart(t, 70),
					newEnginePart(t, valueobject.EngineClassB, 70),
					newShieldPart(t, valueobject.ShieldTypeEnergy),
					newWeaponPart(t, valueobject.WeaponTypeLaser),
				},
			},
			expected: expected{
				wantErr: nil,
			},
		},
		{
			name: "incompatible plasma shield and laser weapon",
			args: args{
				parts: []entity.Part{
					newHullPart(t, 70),
					newEnginePart(t, valueobject.EngineClassB, 70),
					newShieldPart(t, valueobject.ShieldTypePlasma),
					newWeaponPart(t, valueobject.WeaponTypeLaser),
				},
			},
			expected: expected{
				wantErr: pkgerr.FailedPrecondition(errs.ErrIncompatibleParts),
			},
		},
		{
			name: "compatible plasma shield and missile weapon",
			args: args{
				parts: []entity.Part{
					newHullPart(t, 70),
					newEnginePart(t, valueobject.EngineClassB, 70),
					newShieldPart(t, valueobject.ShieldTypePlasma),
					newWeaponPart(t, valueobject.WeaponTypeMissile),
				},
			},
			expected: expected{
				wantErr: nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := compatibilitychecker.New()

			err := svc.Check(tc.args.parts)
			if tc.expected.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tc.expected.wantErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func newHullPart(t *testing.T, strength int32) entity.Part {
	t.Helper()

	properties, err := valueobject.NewHullProperties(strength)
	require.NoError(t, err)

	return newPart(valueobject.PartTypeHull, properties)
}

func newEnginePart(t *testing.T, class valueobject.EngineClass, requiredStrength int32) entity.Part {
	t.Helper()

	properties, err := valueobject.NewEngineProperties(class, requiredStrength)
	require.NoError(t, err)

	return newPart(valueobject.PartTypeEngine, properties)
}

func newShieldPart(t *testing.T, shieldType valueobject.ShieldType) entity.Part {
	t.Helper()

	properties, err := valueobject.NewShieldProperties(shieldType)
	require.NoError(t, err)

	return newPart(valueobject.PartTypeShield, properties)
}

func newWeaponPart(t *testing.T, weaponType valueobject.Type) entity.Part {
	t.Helper()

	properties, err := valueobject.NewWeaponProperties(weaponType)
	require.NoError(t, err)

	return newPart(valueobject.PartTypeWeapon, properties)
}

func newPart(partType valueobject.PartType, properties valueobject.PartProperties) entity.Part {
	return entity.RestorePart(
		gofakeit.UUID(),
		"name",
		"description",
		10,
		10,
		0,
		partType,
		properties,
		time.Now(),
	)
}
