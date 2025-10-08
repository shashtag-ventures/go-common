package gormutil_test

import (
	"context"
	"testing"

	"github.com/shashtag-ventures/go-common/gormutil"
	"github.com/shashtag-ventures/go-common/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type TestModel struct {
	gorm.Model
	Name string
}

func TestGenericRepository(t *testing.T) {
	ctx := context.Background()
	db, teardown := testutil.SetupTestDatabase(ctx)
	defer teardown()

	// We need to migrate the test model
	db.AutoMigrate(&TestModel{})

	repo := gormutil.NewRepository[TestModel](db)

	t.Run("Create and FindByID", func(t *testing.T) {
		testutil.CleanTables(db, "test_models")

		entity := &TestModel{Name: "test"}
		err := repo.Create(ctx, nil, entity)
		require.NoError(t, err)
		assert.NotZero(t, entity.ID)

		found, err := repo.FindByID(ctx, entity.ID)
		require.NoError(t, err)
		assert.Equal(t, "test", found.Name)
	})

	t.Run("FindOneBy", func(t *testing.T) {
		testutil.CleanTables(db, "test_models")

		entity := &TestModel{Name: "findme"}
		err := repo.Create(ctx, nil, entity)
		require.NoError(t, err)

		found, err := repo.FindOneBy(ctx, "name = ?", "findme")
		require.NoError(t, err)
		assert.Equal(t, entity.ID, found.ID)
	})

	t.Run("Update", func(t *testing.T) {
		testutil.CleanTables(db, "test_models")

		entity := &TestModel{Name: "initial"}
		err := repo.Create(ctx, nil, entity)
		require.NoError(t, err)

		entity.Name = "updated"
		err = repo.Update(ctx, nil, entity)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, entity.ID)
		require.NoError(t, err)
		assert.Equal(t, "updated", found.Name)
	})

	t.Run("Delete", func(t *testing.T) {
		testutil.CleanTables(db, "test_models")

		entity := &TestModel{Name: "deleteme"}
		err := repo.Create(ctx, nil, entity)
		require.NoError(t, err)

		err = repo.Delete(ctx, nil, entity)
		require.NoError(t, err)

		_, err = repo.FindByID(ctx, entity.ID)
		assert.Error(t, err)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})

	t.Run("FindAll", func(t *testing.T) {
		testutil.CleanTables(db, "test_models")
		repo.Create(ctx, nil, &TestModel{Name: "one"})
		repo.Create(ctx, nil, &TestModel{Name: "two"})

		all, err := repo.FindAll(ctx)
		require.NoError(t, err)
		assert.Len(t, all, 2)
	})
}
