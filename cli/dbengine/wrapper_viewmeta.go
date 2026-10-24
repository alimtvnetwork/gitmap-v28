package dbengine

import (
	"context"
)

func (w *DBWrapper) recordViewMetaIfPresent(ctx context.Context, name, queryHash, selectSql string) BoolResult {
	if len(queryHash) == 0 {
		return SuccessBool(true)
	}

	saveErr := w.SaveViewMeta(ctx, name, queryHash, selectSql)
	if saveErr != nil {
		return FailureBool(saveErr)
	}

	return SuccessBool(true)
}

func (w *DBWrapper) createAndRegisterView(ctx context.Context, name, selectSql, queryHash string) BoolResult {
	createRes := w.CreateView(ctx, name, selectSql)
	if createRes.IsFailed() {
		return createRes
	}

	return w.recordViewMetaIfPresent(ctx, name, queryHash, selectSql)
}

func (w *DBWrapper) dropAndRecreateView(ctx context.Context, name, selectSql, queryHash string) BoolResult {
	dropRes := w.DropView(ctx, name)
	if dropRes.IsFailed() {
		return dropRes
	}

	return w.createAndRegisterView(ctx, name, selectSql, queryHash)
}

func (w *DBWrapper) recreateAndRegisterView(ctx context.Context, name string, selectSql string, queryHash string) BoolResult {
	valErr := w.ValidateSql(ctx, selectSql)
	if valErr != nil {
		return FailureBool(valErr)
	}

	exists, appErr := w.ViewExists(ctx, name)
	if appErr != nil {
		return FailureBool(appErr)
	}

	if exists {
		return w.dropAndRecreateView(ctx, name, selectSql, queryHash)
	}

	return w.createAndRegisterView(ctx, name, selectSql, queryHash)
}

// CreateViewOrUseViewWithHash checks if a view exists and matches queryHash in __gitmap_view_meta.
// If valid, it reuses the view (0 DDL). If missing or query hash changed, it validates SQL, drops old view, creates updated view, and saves metadata.
func (w *DBWrapper) CreateViewOrUseViewWithHash(ctx context.Context, name string, selectSql string, queryHash string) BoolResult {
	metaErr := w.EnsureViewMetaTable(ctx)
	if metaErr != nil {
		return FailureBool(metaErr)
	}

	exists, appErr := w.ViewExists(ctx, name)
	if appErr != nil {
		return FailureBool(appErr)
	}

	if !exists {
		return w.recreateAndRegisterView(ctx, name, selectSql, queryHash)
	}

	current, currentErr := w.isViewCurrent(ctx, name, queryHash)
	if currentErr != nil {
		return FailureBool(currentErr)
	}

	if current {
		return SuccessBool(true)
	}

	return w.recreateAndRegisterView(ctx, name, selectSql, queryHash)
}

// CreateViewOrUseView inspects if a view exists and contains the required columns or matches query hash.
// If valid, it reuses the view. If missing or schema differs, it validates SQL, drops old view, creates updated view, and saves metadata.
func (w *DBWrapper) CreateViewOrUseView(ctx context.Context, name string, selectSql string, requiredColumns ...string) BoolResult {
	hash := ComputeSqlHash(selectSql)
	if len(requiredColumns) == 0 {
		return w.CreateViewOrUseViewWithHash(ctx, name, selectSql, hash)
	}

	metaErr := w.EnsureViewMetaTable(ctx)
	if metaErr != nil {
		return FailureBool(metaErr)
	}

	exists, appErr := w.ViewExists(ctx, name)
	if appErr != nil {
		return FailureBool(appErr)
	}

	if !exists {
		return w.recreateAndRegisterView(ctx, name, selectSql, hash)
	}

	hasAll, verifyErr := w.verifyViewColumns(ctx, name, requiredColumns)
	if verifyErr != nil {
		return FailureBool(verifyErr)
	}

	if !hasAll {
		return w.recreateAndRegisterView(ctx, name, selectSql, hash)
	}

	saveRes := w.recordViewMetaIfPresent(ctx, name, hash, selectSql)
	if saveRes.IsFailed() {
		return saveRes
	}

	return SuccessBool(true)
}
