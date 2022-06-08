# gqlgen-introspect-filter

[![Go Reference](https://pkg.go.dev/badge/github.com/ec2-software/gqlgen-introspect-filter.svg)](https://pkg.go.dev/github.com/ec2-software/gqlgen-introspect-filter)

Filter GQLGen's Introspection Schema using your application's business logic.

```go

import (
    "github.com/99designs/gqlgen/graphql/handler"
    "yourprojectname/generated/server"
)

// Create the default GQLGen server
exec := handler.NewDefaultServer(
	server.NewExecutableSchema(
        server.Config{Resolvers: resolvers},
	),
)

// Use as a GQLGen plugin
exec.Use(introspectionfilter.Plugin{
	// Write filter functions to choose if various parts are included.
	ReturnField: func(ctx context.Context, fd *ast.FieldDefinition, d *ast.Definition) bool { 
		return fd.Name != "text" 
	},
})
```
