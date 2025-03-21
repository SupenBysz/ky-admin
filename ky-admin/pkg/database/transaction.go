package database

import (
	"context"
	"fmt"

	"github.com/SupenBysz/ky-admin/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// TxFunc 事务处理函数类型
type TxFunc func(tx *gorm.DB) error

// RunInTransaction 在事务中运行函数，自动处理提交和回滚
func RunInTransaction(fn TxFunc) error {
	db := GetDB()
	if db == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return RunInTransactionWithDB(db, fn)
}

// RunInTransactionWithDB 在指定的DB实例上运行事务
func RunInTransactionWithDB(db *gorm.DB, fn TxFunc) error {
	if db == nil {
		return fmt.Errorf("数据库未初始化")
	}

	// 开始事务
	tx := db.Begin()
	if tx.Error != nil {
		logger.Error("开始事务失败", zap.Error(tx.Error))
		return tx.Error
	}

	// 处理panic
	defer func() {
		if r := recover(); r != nil {
			// 回滚事务
			tx.Rollback()
			logger.Error("事务执行时发生panic，已回滚", zap.Any("recover", r))
			// 重新panic
			panic(r)
		}
	}()

	// 执行事务函数
	if err := fn(tx); err != nil {
		// 回滚事务
		tx.Rollback()
		logger.Warn("事务执行失败，已回滚", zap.Error(err))
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		logger.Error("提交事务失败", zap.Error(err))
		return err
	}

	return nil
}

// RunInTransactionWithContext 带有上下文的事务执行函数，支持跟踪ID
func RunInTransactionWithContext(ctx context.Context, fn TxFunc) error {
	db := GetDB()
	if db == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return RunInTransactionWithContextAndDB(ctx, db, fn)
}

// RunInTransactionWithContextAndDB 在指定的DB实例上运行带有上下文的事务
func RunInTransactionWithContextAndDB(ctx context.Context, db *gorm.DB, fn TxFunc) error {
	if db == nil {
		return fmt.Errorf("数据库未初始化")
	}

	traceID := logger.GetTraceID(ctx)
	logFields := []zap.Field{zap.String(logger.TraceIDField, traceID)}

	// 开始事务
	tx := db.WithContext(ctx).Begin()
	if tx.Error != nil {
		logger.ErrorWithContext(ctx, "开始事务失败", zap.Error(tx.Error))
		return tx.Error
	}

	// 处理panic
	defer func() {
		if r := recover(); r != nil {
			// 回滚事务
			tx.Rollback()
			logger.ErrorWithContext(ctx, "事务执行时发生panic，已回滚", append(logFields, zap.Any("recover", r))...)
			// 重新panic
			panic(r)
		}
	}()

	// 执行事务函数
	if err := fn(tx); err != nil {
		// 回滚事务
		tx.Rollback()
		logger.WarnWithContext(ctx, "事务执行失败，已回滚", append(logFields, zap.Error(err))...)
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		logger.ErrorWithContext(ctx, "提交事务失败", append(logFields, zap.Error(err))...)
		return err
	}

	logger.InfoWithContext(ctx, "事务执行成功", logFields...)
	return nil
}

// NestedTransaction 嵌套事务处理
// 如果已经在事务中，则使用SavePoint
// 如果不在事务中，则创建新事务
func NestedTransaction(db *gorm.DB, fn TxFunc) error {
	// 检查是否已经在事务中
	if db.Statement.ConnPool != nil && db.Statement.ConnPool.(*gorm.PreparedStmtDB) != nil {
		// 已经在事务中，使用SavePoint
		sp := fmt.Sprintf("sp_%p", fn)

		// 创建SavePoint
		if err := db.SavePoint(sp).Error; err != nil {
			logger.Error("创建SavePoint失败", zap.Error(err), zap.String("savepoint", sp))
			return err
		}

		// 执行事务函数
		if err := fn(db); err != nil {
			// 回滚到SavePoint
			if rbErr := db.RollbackTo(sp).Error; rbErr != nil {
				logger.Error("回滚SavePoint失败", zap.Error(rbErr), zap.String("savepoint", sp))
				return rbErr
			}
			return err
		}
		return nil
	}

	// 不在事务中，创建新事务
	return RunInTransactionWithDB(db, fn)
}
