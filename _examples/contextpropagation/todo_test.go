package contextpropagation

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ssmccoy/gqlgen/client"
	"github.com/ssmccoy/gqlgen/graphql/handler"
	"github.com/ssmccoy/gqlgen/graphql/handler/extension"
	"github.com/ssmccoy/gqlgen/graphql/handler/transport"
)

func TestTodo(t *testing.T) {
	srv := handler.New(NewExecutableSchema(New()))
	srv.AddTransport(transport.POST{})
	srv.Use(extension.Introspection{})
	c := client.New(srv)

	var resp struct {
		TestTodo struct {
			Text    string
			Context struct {
				Value *string
			}
		}
	}
	c.MustPost(`{ testTodo { text context { value } } }`, &resp)

	require.Equal(t, "Test", resp.TestTodo.Text)
	require.NotNil(t, resp.TestTodo.Context.Value)
	require.Equal(t, "Some value", *resp.TestTodo.Context.Value)
}
