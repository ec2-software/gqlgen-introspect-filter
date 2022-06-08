package introspectionfilter

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/introspection"
	"github.com/vektah/gqlparser/v2/ast"
)

var schema *ast.Schema

type Extension struct {
	ReturnField     FieldFilter
	ReturnType      TypeFilter
	ReturnDirective DirectiveFilter
	ReturnFilter    EnumFilter
}

type FieldFilter func(ctx context.Context, fd *ast.FieldDefinition, d *ast.Definition) bool
type TypeFilter func(ctx context.Context, d *ast.Definition) bool
type DirectiveFilter func(ctx context.Context, dd *ast.DirectiveDefinition) bool
type InputFieldFilter func(ctx context.Context, fd *ast.FieldDefinition) bool
type EnumFilter func(ctx context.Context, ed *ast.EnumValueDefinition) bool

func (e Extension) ExtensionName() string {
	return "IntrospectionFilter"
}
func (e Extension) Validate(s graphql.ExecutableSchema) error {
	schema = s.Schema()
	return nil
}

func (e Extension) InterceptField(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
	res, err = next(ctx)
	if err != nil {
		return
	}

	fc := graphql.GetFieldContext(ctx)

	switch fc.Object {
	case "__Schema":
		switch fc.Field.Name {
		case "types":
			res = e.filterTypes(ctx, res.([]introspection.Type))
		case "directives":
			res = e.filterDirectives(ctx, res.([]introspection.Directive))
		}
	case "__Type":
		iType := fc.Parent.Result.(*introspection.Type)
		tName := iType.Name()
		if tName == nil {
			return
		}

		astType := schema.Types[*tName]
		if astType == nil {
			return
		}

		switch fc.Field.Name {
		case "fields":
			res = e.filterFields(ctx, res.([]introspection.Field), astType)
		case "inputFields":
			res = e.filterInputFields(ctx, res.([]introspection.InputValue), astType)
		case "possibleTypes":
			res = e.filterTypes(ctx, res.([]introspection.Type))
		case "enumValues":
			res = e.filterEnumValues(ctx, res.([]introspection.EnumValue), astType)
		}
	}

	return res, err
}

func (e Extension) filterTypes(ctx context.Context, list []introspection.Type) []introspection.Type {
	if e.ReturnType == nil {
		return list
	}
	fList := make([]introspection.Type, 0, len(list))
	for _, t := range list {
		tName := t.Name()
		if tName != nil {
			astType := schema.Types[*tName]
			if astType == nil {
				continue
			}
			if !e.ReturnType(ctx, astType) {
				continue
			}
		}
		fList = append(fList, t)
	}
	return fList
}

func (e Extension) filterDirectives(ctx context.Context, list []introspection.Directive) []introspection.Directive {
	if e.ReturnDirective == nil {
		return list
	}
	fList := make([]introspection.Directive, 0, len(list))
	for _, x := range list {
		astDirective := schema.Directives[x.Name]
		if astDirective == nil {
			continue
		}
		if !e.ReturnDirective(ctx, astDirective) {
			continue
		}
		fList = append(fList, x)
	}
	return fList
}

func (e Extension) filterFields(ctx context.Context, list []introspection.Field, astType *ast.Definition) []introspection.Field {
	if e.ReturnField == nil {
		return list
	}
	fList := make([]introspection.Field, 0, len(list))
	for _, x := range list {
		astField := astType.Fields.ForName(x.Name)
		if astField == nil {
			continue
		}
		if !e.ReturnField(ctx, astField, astType) {
			continue
		}
		fList = append(fList, x)
	}
	return fList
}

func (e Extension) filterInputFields(ctx context.Context, list []introspection.InputValue, astType *ast.Definition) []introspection.InputValue {
	if e.ReturnField == nil {
		return list
	}
	fList := make([]introspection.InputValue, 0, len(list))
	for _, x := range list {
		astField := astType.Fields.ForName(x.Name)
		if astField == nil {
			continue
		}
		if !e.ReturnField(ctx, astField, astType) {
			continue
		}
		fList = append(fList, x)
	}
	return fList
}

func (e Extension) filterEnumValues(ctx context.Context, list []introspection.EnumValue, astType *ast.Definition) []introspection.EnumValue {
	if e.ReturnFilter == nil {
		return list
	}
	fList := make([]introspection.EnumValue, 0, len(list))
	for _, x := range list {
		astEnum := astType.EnumValues.ForName(x.Name)
		if astEnum == nil {
			continue
		}
		if !e.ReturnFilter(ctx, astEnum) {
			continue
		}
		fList = append(fList, x)
	}
	return fList
}
