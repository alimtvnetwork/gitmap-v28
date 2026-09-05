package dbengine

import (
	"context"
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

// Query returns a new fluent QueryBuilder for this repository.
func (r *Repository[T, F]) Query() *QueryBuilder[T, F] {
	return NewQueryBuilder(r)
}

// First queries a single record matching one field value, returning an EntityResult envelope.
func (r *Repository[T, F]) First(ctx context.Context, field F, value any) EntityResult[T] {
	return r.Query().WhereEq(field, value).First(ctx)
}

// FindById queries a single record by its uint64 primary key.
func (r *Repository[T, F]) FindById(ctx context.Context, idField F, id uint64) EntityResult[T] {
	return r.First(ctx, idField, id)
}

// FindBy queries up to limit records matching one field value, returning a ListResult envelope.
func (r *Repository[T, F]) FindBy(ctx context.Context, field F, value any, limit int) ListResult[T] {
	return r.Query().WhereEq(field, value).Limit(limit).FindAll(ctx)
}

// FindBy2 queries up to limit records matching two field values, returning a ListResult envelope.
func (r *Repository[T, F]) FindBy2(ctx context.Context, field1 F, val1 any, field2 F, val2 any, limit int) ListResult[T] {
	return r.Query().WhereEq(field1, val1).WhereEq(field2, val2).Limit(limit).FindAll(ctx)
}

// FindAll queries up to limit records from the table, returning a ListResult envelope.
func (r *Repository[T, F]) FindAll(ctx context.Context, limit int) ListResult[T] {
	return r.Query().Limit(limit).FindAll(ctx)
}

// Count queries the count of records matching a field value, returning an Int64Result envelope.
func (r *Repository[T, F]) Count(ctx context.Context, field F, value any) Int64Result {
	return r.Query().WhereEq(field, value).Count(ctx)
}

// CountAll queries the total count of records in the table, returning an Int64Result envelope.
func (r *Repository[T, F]) CountAll(ctx context.Context) Int64Result {
	return r.Query().Count(ctx)
}

// DeleteBy deletes records matching a field value, returning a RowsAffectedResult envelope.
func (r *Repository[T, F]) DeleteBy(ctx context.Context, field F, value any) RowsAffectedResult {
	return r.Query().WhereEq(field, value).Delete(ctx)
}
