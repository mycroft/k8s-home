package kubehelpers

import "context"

type ContextValues struct {
	Debug bool
}

type key int

const (
	ValueKey key = iota
)

// ContextGetDebug returns true if the passed context has a debug flag to true, else false
func ContextGetDebug(ctx context.Context) bool {
	values := ctx.Value(ValueKey)
	if values == nil {
		return false
	}

	return values.(ContextValues).Debug
}
