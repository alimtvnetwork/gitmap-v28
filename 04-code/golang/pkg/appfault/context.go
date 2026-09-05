package appfault

// WithContext returns a new immutable AppError with the attached key-value context.
// The receiver is never mutated in-place, guaranteeing strict immutability.
func (e *AppError) WithContext(key string, value any) *AppError {
	if e == nil {
		return nil
	}

	cloned := e.clone()
	cloned.ctx.Set(key, value)

	return cloned
}

// WithOp returns a new immutable AppError with the operation context.
func (e *AppError) WithOp(op string) *AppError {
	return e.WithContext("Op", op)
}

// WithSeverity returns a new immutable AppError with the severity level.
func (e *AppError) WithSeverity(severity SeverityType) *AppError {
	return e.WithContext("Severity", severity.Name())
}

// WithPriority returns a new immutable AppError with the priority level.
func (e *AppError) WithPriority(priority PriorityType) *AppError {
	return e.WithContext("Priority", priority.Name())
}

// WithUrl returns a new immutable AppError with the request URL context.
func (e *AppError) WithUrl(url string) *AppError {
	return e.WithContext("Url", url)
}

// WithStatusCode returns a new immutable AppError with the specified HTTP status code.
func (e *AppError) WithStatusCode(statusCode int) *AppError {
	if e == nil {
		return nil
	}

	cloned := e.clone()
	cloned.statusCode = statusCode
	cloned.ctx.Set("StatusCode", statusCode)

	return cloned
}

// WithEndpoint returns a new immutable AppError with the API endpoint context.
func (e *AppError) WithEndpoint(endpoint string) *AppError {
	return e.WithContext("Endpoint", endpoint)
}

// WithSiteId returns a new immutable AppError with the target Site ID context.
func (e *AppError) WithSiteId(siteId int64) *AppError {
	return e.WithContext("SiteId", siteId)
}

// WithSnapshotId returns a new immutable AppError with the snapshot identifier.
func (e *AppError) WithSnapshotId(snapshotId string) *AppError {
	return e.WithContext("SnapshotId", snapshotId)
}

// WithSlug returns a new immutable AppError with the entity slug.
func (e *AppError) WithSlug(slug string) *AppError {
	return e.WithContext("Slug", slug)
}

// WithPluginContext returns a new immutable AppError with plugin ID and slug.
func (e *AppError) WithPluginContext(pluginId int64, slug string) *AppError {
	return e.WithContext("PluginId", pluginId).WithSlug(slug)
}

// WithCaller returns a new immutable AppError with the specified caller site metadata.
func (e *AppError) WithCaller(caller CallerInfo) *AppError {
	if e == nil {
		return nil
	}

	cloned := e.clone()
	cloned.caller = caller

	return cloned
}

// Context returns a copy of the underlying diagnostic metadata ContextMap.
func (e *AppError) Context() ContextMap {
	if e == nil || e.ctx == nil {
		return NewContextMap()
	}

	return e.ctx.Clone()
}
