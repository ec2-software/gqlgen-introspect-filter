package introspectionfilter

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/introspection"
	"github.com/vektah/gqlparser/v2/ast"
)

type Plugin struct {
	schema *ast.Schema

	FieldFilter     FieldFilter
	TypeFilter      TypeFilter
	DirectiveFilter DirectiveFilter
	EnumFilter      EnumFilter
}

type FieldFilter func(ctx context.Context, d *ast.FieldDefinition) bool
type TypeFilter func(ctx context.Context, d *ast.Definition) bool
type DirectiveFilter func(ctx context.Context, d *ast.DirectiveDefinition) bool
type InputFieldFilter func(ctx context.Context, d *ast.FieldDefinition) bool
type EnumFilter func(ctx context.Context, d *ast.EnumValueDefinition) bool

func (*Plugin) ExtensionName() string {
	return "IntrospectionFilter"
}
func (p *Plugin) Validate(schema graphql.ExecutableSchema) error {
	p.schema = schema.Schema()
	return nil
}

func (p *Plugin) InterceptField(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
	res, err = next(ctx)
	if err != nil {
		return
	}

	fc := graphql.GetFieldContext(ctx)
	schema := p.schema

	switch fc.Object {
	case "__Schema":
		switch fc.Field.Name {
		case "types":
			res = p.filterTypes(ctx, res.([]introspection.Type))
		case "directives":
			res = p.filterDirectives(ctx, res.([]introspection.Directive))
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
			res = p.filterFields(ctx, res.([]introspection.Field), astType)
		case "inputFields":
			res = p.filterInputFields(ctx, res.([]introspection.InputValue), astType)
		case "possibleTypes":
			res = p.filterTypes(ctx, res.([]introspection.Type))
		case "enumValues":
			res = p.filterEnumValues(ctx, res.([]introspection.EnumValue), astType)
		}
	}

	return res, err
}

func (filter *Plugin) filterTypes(ctx context.Context, list []introspection.Type) []introspection.Type {
	if filter.TypeFilter == nil {
		return list
	}
	fList := make([]introspection.Type, 0, len(list))
	for _, t := range list {
		tName := t.Name()
		if tName != nil {
			astType := filter.schema.Types[*tName]
			if astType == nil {
				continue
			}
			if !filter.TypeFilter(ctx, astType) {
				continue
			}
		}
		fList = append(fList, t)
	}
	return fList
}

func (p *Plugin) filterDirectives(ctx context.Context, list []introspection.Directive) []introspection.Directive {
	if p.DirectiveFilter == nil {
		return list
	}
	fList := make([]introspection.Directive, 0, len(list))
	for _, x := range list {
		astDirective := p.schema.Directives[x.Name]
		if astDirective == nil {
			continue
		}
		if !p.DirectiveFilter(ctx, astDirective) {
			continue
		}
		fList = append(fList, x)
	}
	return fList
}

func (p *Plugin) filterFields(ctx context.Context, list []introspection.Field, astType *ast.Definition) []introspection.Field {
	if p.FieldFilter == nil {
		return list
	}
	fList := make([]introspection.Field, 0, len(list))
	for _, x := range list {
		astField := astType.Fields.ForName(x.Name)
		if astField == nil {
			continue
		}
		if !p.FieldFilter(ctx, astField) {
			continue
		}
		fList = append(fList, x)
	}
	return fList
}

func (p *Plugin) filterInputFields(ctx context.Context, list []introspection.InputValue, astType *ast.Definition) []introspection.InputValue {
	if p.FieldFilter == nil {
		return list
	}
	fList := make([]introspection.InputValue, 0, len(list))
	for _, x := range list {
		astField := astType.Fields.ForName(x.Name)
		if astField == nil {
			continue
		}
		if !p.FieldFilter(ctx, astField) {
			continue
		}
		fList = append(fList, x)
	}
	return fList
}

func (p *Plugin) filterEnumValues(ctx context.Context, list []introspection.EnumValue, astType *ast.Definition) []introspection.EnumValue {
	if p.EnumFilter == nil {
		return list
	}
	fList := make([]introspection.EnumValue, 0, len(list))
	for _, x := range list {
		astEnum := astType.EnumValues.ForName(x.Name)
		if astEnum == nil {
			continue
		}
		if !p.EnumFilter(ctx, astEnum) {
			continue
		}
		fList = append(fList, x)
	}
	return fList
}
