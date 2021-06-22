package introspectionfilter

import (
	"context"
	"sort"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/introspection"
)

type SortPlugin struct{}

func (SortPlugin) ExtensionName() string {
	return "IntrospectionSort"
}
func (SortPlugin) Validate(schema graphql.ExecutableSchema) error {
	return nil
}

func (SortPlugin) InterceptField(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
	res, err = next(ctx)
	if err != nil {
		return
	}

	fc := graphql.GetFieldContext(ctx)

	// __ is reserved for introspection types
	if strings.HasPrefix(fc.Object, "__") {
		switch x := res.(type) {
		case []introspection.Type:
			sort.SliceStable(x, func(i, j int) bool {
				in, jn := x[i].Name(), x[j].Name()
				if in == nil {
					if jn == nil {
						return x[i].Description() < x[j].Description()
					}
					return true
				}
				if jn == nil {
					return false
				}
				return *in < *jn
			})
		case []introspection.Directive:
			sort.SliceStable(x, func(i, j int) bool { return x[i].Name < x[j].Name })
		case []introspection.EnumValue:
			sort.SliceStable(x, func(i, j int) bool { return x[i].Name < x[j].Name })
		case []introspection.Field:
			sort.SliceStable(x, func(i, j int) bool { return x[i].Name < x[j].Name })
		case []introspection.InputValue:
			sort.SliceStable(x, func(i, j int) bool { return x[i].Name < x[j].Name })
		}
	}

	return
}
