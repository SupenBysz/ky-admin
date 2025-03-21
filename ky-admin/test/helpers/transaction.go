package helpers

import (
	"gorm.io/gorm"
)

// WithTransaction 在事务中执行操作
func WithTransaction(db *gorm.DB, fn func(tx interface{}) error) error {
	return db.Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}
