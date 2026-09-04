package appfault

// WithContext attaches a generic key-value pair to the AppError context.
func (e *AppError) WithContext(key string, value any) *AppError {
	if e == nil {
		return nil
	}

	if e.ctx == nil {
		e.ctx = NewContextMap()
	}

	e.ctx.Set(key, value)

	return e
}

// WithOp attaches an operation name to the diagnostic context.
func (e *AppError) WithOp(op string) *AppError {
	return e.WithContext("Op", op)
}

// WithSeverity attaches a typed severity level to the diagnostic context.
func (e *AppError) WithSeverity(severity SeverityType) *AppError {
	return e.WithContext("Severity", severity.Name())
}

// WithPriority attaches a typed priority level to the context.
func (e *AppError) WithPriority(priority PriorityType) *AppError {
	return e.WithContext("Priority", priority.Name())
}

// WithUrl attaches a request URL to the diagnostic context.
func (e *AppError) WithUrl(url string) *AppError {
	return e.WithContext("Url", url)
}

// WithStatusCode attaches an HTTP status code to the error.
func (e *AppError) WithStatusCode(statusCode int) *AppError {
	if e != nil {
		e.statusCode = statusCode
	}

	return e.WithContext("StatusCode", statusCode)
}

// WithEndpoint attaches an API endpoint path to the context.
func (e *AppError) WithEndpoint(endpoint string) *AppError {
	return e.WithContext("Endpoint", endpoint)
}

// WithSiteId attaches a target Site ID to the context.
func (e *AppError) WithSiteId(siteId int64) *AppError {
	return e.WithContext("SiteId", siteId)
}

// WithSnapshotId attaches a snapshot identifier to the context.
func (e *AppError) WithSnapshotId(snapshotId string) *AppError {
	return e.WithContext("SnapshotId", snapshotId)
}

// WithSlug attaches a plugin or entity slug to the context.
func (e *AppError) WithSlug(slug string) *AppError {
	return e.WithContext("Slug", slug)
}

// WithPluginContext attaches both plugin ID and slug simultaneously.
func (e *AppError) WithPluginContext(pluginId int64, slug string) *AppError {
	return e.WithContext("PluginId", pluginId).WithSlug(slug)
}

// Context returns a copy of the underlying diagnostic metadata ContextMap.
func (e *AppError) Context() ContextMap {
	if e == nil || e.ctx == nil {
		return NewContextMap()
	}

	return e.ctx.Clone()
}
