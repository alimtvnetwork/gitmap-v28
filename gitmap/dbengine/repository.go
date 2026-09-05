package dbengine

import (
	"context"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// RowScanner abstracts single-row and multi-row scanners.
type RowScanner interface {
	Scan(dest ...any) error
}

// ModelScanner defines a mapping function from a database row to a typed model pointer.
type ModelScanner[T any] func(s RowScanner) (*T, error)

// Repository provides a type-safe generic data access repository.
type Repository[T any, F ~string] struct {
	db        *DbWrapper
	tableName string
	scanner   ModelScanner[T]
}

// NewRepository initializes a Repository for a table and model scanner.
func NewRepository[T any, F ~string](db *DbWrapper, tableName string, scanner ModelScanner[T]) *Repository[T, F] {
	return &Repository[T, F]{
		db:        db,
		tableName: tableName,
		scanner:   scanner,
	}
}

// First queries a single record matching one field value.
func (r *Repository[T, F]) First(ctx context.Context, field F, value any) (*T, *apperror.AppError) {
	query := r.db.compiler.CompileSearch(r.tableName, []string{string(field)}, 1)
	row, appErr := r.db.QueryRow(ctx, query, value)
	if appErr != nil {
		return nil, appErr
	}

	item, scanErr := r.scanner(row)
	if scanErr != nil {
		return nil, apperror.WrapSimple(scanErr, fmt.Sprintf("scan single %s", r.tableName))
	}
	return item, nil
}

// FindBy queries up to limit records matching one field value.
func (r *Repository[T, F]) FindBy(ctx context.Context, field F, value any, limit int) ([]T, *apperror.AppError) {
	query := r.db.compiler.CompileSearch(r.tableName, []string{string(field)}, limit)
	rows, appErr := r.db.Query(ctx, query, value)
	if appErr != nil {
		return nil, appErr
	}
	defer rows.Close()

	var items []T
	for rows.Next() {
		item, scanErr := r.scanner(rows)
		if scanErr != nil {
			return nil, apperror.WrapSimple(scanErr, fmt.Sprintf("scan %s row", r.tableName))
		}
		items = append(items, *item)
	}
	return items, nil
}

// FindBy2 queries up to limit records matching two field values.
func (r *Repository[T, F]) FindBy2(ctx context.Context, field1 F, val1 any, field2 F, val2 any, limit int) ([]T, *apperror.AppError) {
	query := r.db.compiler.CompileSearch(r.tableName, []string{string(field1), string(field2)}, limit)
	rows, appErr := r.db.Query(ctx, query, val1, val2)
	if appErr != nil {
		return nil, appErr
	}
	defer rows.Close()

	var items []T
	for rows.Next() {
		item, scanErr := r.scanner(rows)
		if scanErr != nil {
			return nil, apperror.WrapSimple(scanErr, fmt.Sprintf("scan %s row with 2 params", r.tableName))
		}
		items = append(items, *item)
	}
	return items, nil
}

// FindAll queries up to limit records from the table without filters.
func (r *Repository[T, F]) FindAll(ctx context.Context, limit int) ([]T, *apperror.AppError) {
	query := r.db.compiler.CompileSearch(r.tableName, []string{}, limit)
	rows, appErr := r.db.Query(ctx, query)
	if appErr != nil {
		return nil, appErr
	}
	defer rows.Close()

	var items []T
	for rows.Next() {
		item, scanErr := r.scanner(rows)
		if scanErr != nil {
			return nil, apperror.WrapSimple(scanErr, fmt.Sprintf("scan all %s", r.tableName))
		}
		items = append(items, *item)
	}
	return items, nil
}
