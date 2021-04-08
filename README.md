# gqlgen-introspect-filter
Filter GQLGen's Introspection by your application's business logic.

```go
// Use as a GQLGen plugin
exec.Use(introspectionfilter.Plugin{
	Schema:      schema,
	FieldFilter: func(fd *ast.FieldDefinition) bool { return fd.Name != "text" }, // Write filter functions to choose if various parts are included.
})
```
