package testutil

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupTestDatabase(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	ctx := context.Background()
	db, teardown := SetupTestDatabase(ctx)
	defer teardown()

	assert.NotNil(t, db)

	// Test table creation and cleaning
	type TestModel struct {
		ID   uint   `gorm:"primaryKey"`
		Name string
	}

	err := db.AutoMigrate(&TestModel{})
	assert.NoError(t, err)

	db.Create(&TestModel{Name: "test"})
	
	var count int64
	db.Model(&TestModel{}).Count(&count)
	assert.Equal(t, int64(1), count)

	CleanTables(db, "test_models")

	db.Model(&TestModel{}).Count(&count)
	assert.Equal(t, int64(0), count)
}
