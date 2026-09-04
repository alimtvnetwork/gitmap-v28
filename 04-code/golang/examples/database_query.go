package examples

import (
	"context"
	"database/sql"
	"errors"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
	"coding-guidelines/common/pkg/result"
)

// PluginSummary represents an active plugin record in the system.
type PluginSummary struct {
	Id       int64  `json:"id"`
	Slug     string `json:"slug"`
	Name     string `json:"name"`
	IsActive bool   `json:"isActive"`
}

// PluginRepository provides data access methods using Result wrappers.
type PluginRepository struct {
	db *sql.DB
}

// NewPluginRepository constructs a repository instance.
func NewPluginRepository(db *sql.DB) *PluginRepository {
	return &PluginRepository{db: db}
}

// checkSpecialPluginIds simulates database failure states for test cases.
func checkSpecialPluginIds(id int64) *appfault.AppError {
	if id == 404 {
		return appfault.New(errtype.NotFound, "plugin record not found").WithOp("repo.find")
	}

	if id == 500 {
		return appfault.Wrap(errtype.Database, errors.New("connection reset by peer"), "connection failed").WithOp("repo.find")
	}

	return nil
}

// validatePluginId checks input and simulates database error states.
func validatePluginId(id int64) *appfault.AppError {
	if id <= 0 {
		return appfault.New(errtype.Validation, "plugin id must be positive")
	}

	return checkSpecialPluginIds(id)
}

// FindById queries a single plugin record by Id.
func (r *PluginRepository) FindById(ctx context.Context, id int64) result.Result[PluginSummary] {
	if err := validatePluginId(id); err != nil {
		return result.FailureResult[PluginSummary](err)
	}

	return result.SuccessResult(PluginSummary{
		Id:       id,
		Slug:     "seo-optimizer",
		Name:     "SEO Optimizer Pro",
		IsActive: true,
	})
}

// ListActive retrieves all active plugins as a ResultSlice.
func (r *PluginRepository) ListActive(ctx context.Context) appfault.ResultSlice[PluginSummary] {
	items := []PluginSummary{
		{Id: 1, Slug: "cache-booster", Name: "Cache Booster", IsActive: true},
		{Id: 2, Slug: "security-shield", Name: "Security Shield", IsActive: true},
	}

	return appfault.OkSlice(items)
}
