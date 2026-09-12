package templates

const (
	kindIgnore             = "ignore"
	kindAttributes         = "attributes"
	kindLFS                = "lfs"
	userTemplatesDirName   = ".gitmap"
	userTemplatesSubdir    = "templates"
	embedAssetsRoot        = "assets"
	templateExtIgnore      = ".gitignore"
	templateExtAttributes  = ".gitattributes"
	templateHeaderSource   = "# source:"
	templateHeaderKind     = "# kind:"
	templateHeaderLang     = "# lang:"
	templateHeaderVersion  = "# version:"
	errTemplateNotFound    = "template not found: kind=%s lang=%s"
	errTemplateMaterialize = "failed to materialize templates to %s: %w"
	errTemplateUserDir     = "failed to resolve user templates dir: %w"
	errTemplateRead        = "failed to read template %s: %w"
)
