package gormutil

import (
	"context"

	"gorm.io/gorm"
)

type txKey struct{}

// WithTransaction returns a context with the given transaction.
func WithTransaction(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// GetTransaction retrieves the transaction from the context if it exists.
func GetTransaction(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return nil
}
