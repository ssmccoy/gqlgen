package graphql

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser/v2/ast"
)

type (
	testResolverRoot   struct{}
	testDirectiveRoot  struct{}
	testComplexityRoot struct{}
)

func newTestExecutionContextState(
	opCtx *OperationContext,
	schemaData *ast.Schema,
	parsedSchema *ast.Schema,
) *ExecutionContextState[testResolverRoot, testDirectiveRoot, testComplexityRoot] {
	if opCtx == nil {
		opCtx = &OperationContext{}
	}
	return NewExecutionContextState[testResolverRoot, testDirectiveRoot, testComplexityRoot](
		opCtx,
		&ExecutableSchemaState[testResolverRoot, testDirectiveRoot, testComplexityRoot]{
			SchemaData: schemaData,
		},
		parsedSchema,
	)
}

func makeSchemaWithType(typeName string) *ast.Schema {
	query := &ast.Definition{Name: "Query", Kind: ast.Object}
	typ := &ast.Definition{Name: typeName, Kind: ast.Object}

	return &ast.Schema{
		Query: query,
		Types: map[string]*ast.Definition{
			"Query":  query,
			typeName: typ,
		},
		Directives: map[string]*ast.DirectiveDefinition{},
	}
}

func TestExecutionContextState_Schema(t *testing.T) {
	schemaData := makeSchemaWithType("SchemaDataType")
	parsedSchema := makeSchemaWithType("ParsedType")

	ec := newTestExecutionContextState(
		&OperationContext{},
		schemaData,
		parsedSchema,
	)

	assert.Same(t, schemaData, ec.Schema())
}

func TestExecutionContextState_Schema_FallsBackToParsedSchema(t *testing.T) {
	parsedSchema := makeSchemaWithType("ParsedType")

	ec := newTestExecutionContextState(
		&OperationContext{},
		nil,
		parsedSchema,
	)

	assert.Same(t, parsedSchema, ec.Schema())
}

func TestExecutionContextState_IntrospectionDisabled(t *testing.T) {
	ec := newTestExecutionContextState(
		&OperationContext{DisableIntrospection: true},
		nil,
		makeSchemaWithType("Foo"),
	)

	schema, schemaErr := ec.IntrospectSchema()
	require.Error(t, schemaErr)
	require.EqualError(t, schemaErr, "introspection disabled")
	assert.Nil(t, schema)

	typ, typeErr := ec.IntrospectType("Foo")
	require.Error(t, typeErr)
	require.EqualError(t, typeErr, "introspection disabled")
	assert.Nil(t, typ)
}

func TestExecutionContextState_IntrospectType(t *testing.T) {
	ec := newTestExecutionContextState(
		&OperationContext{},
		nil,
		makeSchemaWithType("Foo"),
	)

	typ, err := ec.IntrospectType("Foo")
	require.NoError(t, err)
	require.NotNil(t, typ)

	name := typ.Name()
	require.NotNil(t, name)
	assert.Equal(t, "Foo", *name)

	missing, missingErr := ec.IntrospectType("Missing")
	require.NoError(t, missingErr)
	assert.Nil(t, missing)
}

func TestExecutionContextState_ProcessDeferredGroup_IncrementsPendingAndPropagates(t *testing.T) {
	ec := newTestExecutionContextState(
		&OperationContext{},
		nil,
		makeSchemaWithType("Foo"),
	)

	ctx := WithResponseContext(context.Background(), DefaultErrorPresenter, DefaultRecover)
	fieldSet := NewFieldSet([]CollectedField{{Field: &ast.Field{Alias: "value"}}})
	fieldSet.Concurrently(0, func(ctx context.Context) Marshaler {
		return MarshalString("ok")
	})

	path := ast.Path{ast.PathName("query")}
	label := "group-1"

	ec.ProcessDeferredGroup(DeferredGroup{
		Path:     path,
		Label:    label,
		FieldSet: fieldSet,
		Context:  ctx,
	})

	require.Len(t, ec.DeferredResults, 1)
	result := ec.DeferredResults[0]
	assert.Equal(t, path, result.Path)
	assert.Equal(t, label, result.Label)
	assert.Same(t, fieldSet, result.Result)
	assert.Nil(t, result.Errors)
}

func TestExecutionContextState_ProcessDeferredGroup_NullsOnInvalidAndIsolatesErrors(t *testing.T) {
	ec := newTestExecutionContextState(
		&OperationContext{},
		nil,
		makeSchemaWithType("Foo"),
	)

	ctx := WithResponseContext(context.Background(), DefaultErrorPresenter, DefaultRecover)
	AddError(ctx, errors.New("parent error"))

	fieldSet := NewFieldSet([]CollectedField{{Field: &ast.Field{Alias: "value"}}})
	fieldSet.Concurrently(0, func(ctx context.Context) Marshaler {
		AddError(ctx, errors.New("deferred error"))
		fieldSet.Invalids++
		return MarshalString("ignored")
	})

	ec.ProcessDeferredGroup(DeferredGroup{
		Path:     ast.Path{ast.PathName("query")},
		Label:    "group-2",
		FieldSet: fieldSet,
		Context:  ctx,
	})

	require.Len(t, ec.DeferredResults, 1)
	result := ec.DeferredResults[0]
	assert.Same(t, Null, result.Result)

	require.Len(t, result.Errors, 1)
	assert.Equal(t, "deferred error", result.Errors[0].Message)

	parentErrors := GetErrors(ctx)
	require.Len(t, parentErrors, 1)
	assert.Equal(t, "parent error", parentErrors[0].Message)
}
