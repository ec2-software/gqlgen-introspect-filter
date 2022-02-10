package introspectionfilter_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/executor"
	"github.com/99designs/gqlgen/graphql/introspection"
	"github.com/bradleyjkemp/cupaloy"
	introspectionfilter "github.com/ec2-software/gqlgen-introspect-filter"
	"github.com/ec2-software/gqlgen-introspect-filter/internal/chat"
	"github.com/vektah/gqlparser/v2/ast"
)

func TestPlugin(t *testing.T) {
	xs := chat.NewExecutableSchema(chat.New())

	exec := executor.New(xs)
	exec.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		graphql.GetOperationContext(ctx).DisableIntrospection = false
		return next(ctx)
	})
	exec.Use(introspectionfilter.SortPlugin{})
	exec.Use(&introspectionfilter.Plugin{
		FieldFilter: func(ctx context.Context, fd *ast.FieldDefinition) bool { return fd.Name != "text" },
	})
	ctx := context.Background()

	rc, errl := exec.CreateOperationContext(graphql.StartOperationTrace(ctx), &graphql.RawParams{
		Query:    introspection.Query,
		ReadTime: graphql.TraceTiming{Start: time.Now()},
	})
	if len(errl) > 0 {
		t.Fatal(errl)
	}
	handler, ctx := exec.DispatchOperation(ctx, rc)
	response := handler(ctx)
	if len(response.Errors) > 0 {
		t.Fatal(response.Errors)
	}

	snapshotter := cupaloy.New(cupaloy.SnapshotSubdirectory("testdata"), cupaloy.SnapshotFileExtension(".json"))
	snapshotter.SnapshotT(t, normalizeJSON(response.Data))
}

func normalizeJSON(in []byte) []byte {
	var d map[string]interface{}
	err := json.Unmarshal(in, &d)
	if err != nil {
		panic(err)
	}
	b, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		panic(err)
	}
	return b
}
