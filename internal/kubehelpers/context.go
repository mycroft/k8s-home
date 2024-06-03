package kubehelpers

import "context"

type ContextValues struct {
	Debug bool
}

// ContextGetDebug returns true if the passed context has a debug flag to true, else false
func ContextGetDebug(ctx context.Context) bool {
	values := ctx.Value("values")
	if values == nil {
		return false
	}

	return values.(ContextValues).Debug
}
