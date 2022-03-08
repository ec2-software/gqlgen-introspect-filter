package introspectionfilter

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
)

// SortPlugin is a plugin that sorts the result of introspection queries.
// Deprecated: GQLGen now sorts the result of introspection queries automatically. This plugin is now a no-op.
type SortPlugin struct{}

func (SortPlugin) ExtensionName() string {
	return "IntrospectionSort"
}
func (SortPlugin) Validate(schema graphql.ExecutableSchema) error {
	return nil
}

func (SortPlugin) InterceptField(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
	return next(ctx)
}
