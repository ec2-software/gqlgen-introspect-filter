package introspectionfilter_test

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/99designs/gqlgen/example/chat"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/executor"
	"github.com/99designs/gqlgen/graphql/introspection"
	introspectionfilter "github.com/ec2-software/gqlgen-introspect-filter"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
)

//go:embed expected_result.json
var expectedResult []byte

func init() {
	expectedResult = normalizeJSON(expectedResult)
}

func TestPlugin(t *testing.T) {
	xs := chat.NewExecutableSchema(chat.New())

	exec := executor.New(xs)
	exec.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		graphql.GetOperationContext(ctx).DisableIntrospection = false
		return next(ctx)
	})
	exec.Use(introspectionfilter.SortPlugin{})
	exec.Use(introspectionfilter.Plugin{
		Schema:      xs.Schema(),
		FieldFilter: func(fd *ast.FieldDefinition) bool { return fd.Name != "text" },
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

	data := normalizeJSON(response.Data)
	if !bytes.Equal(expectedResult, data) {
		var aJson map[string]interface{}

		differ := gojsondiff.New()
		diff, err := differ.Compare(expectedResult, data)
		if err != nil {
			t.Fatal(err)
		}
		err = json.Unmarshal(expectedResult, &aJson)
		if err != nil {
			t.Fatal(err)
		}

		// Write the results for easy comparison
		_ = os.WriteFile("result.json", data, os.ModePerm)

		f := formatter.NewAsciiFormatter(aJson, formatter.AsciiFormatterDefaultConfig)
		res, err := f.Format(diff)
		if err != nil {
			t.Fatal(err)
		}
		t.Error(res)
	}
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
