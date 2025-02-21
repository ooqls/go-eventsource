package eventsourcingv1

import (
	"github.com/jmoiron/sqlx"
	"github.com/ooqls/go-log"
	"go.uber.org/zap"
)

var l *zap.Logger = log.NewLogger("sql")

type Iterator[T any] func() *T

func NewArrayIterator[T any](arr []T) Iterator[T] {
	i := 0
	return func() *T {
		if i < len(arr) {
			v := arr[i]
			i++
			return &v
		}
		return nil
	}
}

func NewSQLIterator[T any](rows *sqlx.Rows) Iterator[T] {
	return func() *T {
		var v T
		if rows.Next() {
			if err := rows.StructScan(v); err != nil {
				l.Error("Failed to scan row", zap.Error(err))
				return nil
			}
			return &v
		} else {
			err := rows.Close()
			if err != nil {
				l.Error("Error closing rows", zap.Error(err))
			}
		}

		if err := rows.Err(); err != nil {
			l.Error("Error iterating rows", zap.Error(err))
		}

		return nil
	}
}
