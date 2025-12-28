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

	t.Run("Find", func(t *testing.T) {
		testutil.CleanTables(db, "test_models")
		repo.Create(ctx, nil, &TestModel{Name: "search1"})
		repo.Create(ctx, nil, &TestModel{Name: "search2"})
		repo.Create(ctx, nil, &TestModel{Name: "other"})

		results, err := repo.Find(ctx, "name LIKE ?", "search%")
		require.NoError(t, err)
		assert.Len(t, results, 2)
	})

	t.Run("DB and getDB", func(t *testing.T) {
		assert.NotNil(t, repo.DB(ctx))
		
		// Indirectly test getDB via Create with transaction
		tx := db.Begin()
		err := repo.Create(ctx, tx, &TestModel{Name: "tx-test"})
		assert.NoError(t, err)
		tx.Rollback()
		
		// Verify not created due to rollback
		_, err = repo.FindOneBy(ctx, "name = ?", "tx-test")
		assert.Error(t, err)
	})

	t.Run("FindAll Error", func(t *testing.T) {
		// Create a repository with a closed DB to force an error if possible, 
		// or just a repo for a non-existent table.
		type NonExistent struct{ gorm.Model }
		repoErr := gormutil.NewRepository[NonExistent](db)
		_, err := repoErr.FindAll(ctx)
		assert.Error(t, err)
	})

	t.Run("Find Error", func(t *testing.T) {
		type NonExistent struct{ gorm.Model }
		repoErr := gormutil.NewRepository[NonExistent](db)
		_, err := repoErr.Find(ctx, "id = ?", 1)
		assert.Error(t, err)
	})

	t.Run("Failures", func(t *testing.T) {
		testutil.CleanTables(db, "test_models")

		// FindByID failure
		_, err := repo.FindByID(ctx, 9999)
		assert.Error(t, err)

		// FindOneBy failure
		_, err = repo.FindOneBy(ctx, "name = ?", "none")
		assert.Error(t, err)

		// Create failure (nil entity)
		err = repo.Create(ctx, nil, nil)
		assert.Error(t, err)

		// Update failure (nil entity)
		err = repo.Update(ctx, nil, nil)
		assert.Error(t, err)

		// Delete non-existent
		err = repo.Delete(ctx, nil, &TestModel{Model: gorm.Model{ID: 9999}})
		assert.NoError(t, err)
	})
}
