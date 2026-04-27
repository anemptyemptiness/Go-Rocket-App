package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/model"
	inventoryrepo "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/repository/part"
)

func TestListParts_Success_FilterByType(t *testing.T) {
	t.Parallel()

	var (
		ctx          = context.Background()
		partTypeHull = model.PartTypeHull
	)

	repo := inventoryrepo.New()

	parts, err := repo.ListParts(ctx, nil, partTypeHull)
	require.NoError(t, err)
	require.NotEmpty(t, parts)

	for _, part := range parts {
		require.Equal(t, part.PartType, partTypeHull)
	}
}

func TestListParts_Success_PartTypeUnspecified(t *testing.T) {
	t.Parallel()

	var (
		ctx                 = context.Background()
		partTypeUnspecified = model.PartTypeUnspecified
	)

	repo := inventoryrepo.New()

	parts, err := repo.ListParts(ctx, nil, partTypeUnspecified)
	require.NoError(t, err)
	require.NotEmpty(t, parts)
	require.Len(t, parts, 7)
}
