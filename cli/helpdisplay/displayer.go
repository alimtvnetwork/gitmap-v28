// Package helpdisplay provides a DRY, themeable object model for rendering
// CLI help output: small interfaces, concrete structs, bound via constructors.
package helpdisplay

// Displayer is the interface every help display implements.
type Displayer interface {
	// Render builds the full help text for the given context.
	Render(ctx RenderContext) string
	// Print writes the rendered help text to stdout.
	Print(ctx RenderContext)
}
