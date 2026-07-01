package singlefile

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ssmccoy/gqlgen/client"
	"github.com/ssmccoy/gqlgen/graphql/handler"
	"github.com/ssmccoy/gqlgen/graphql/handler/transport"
)

func TestTypeFallback(t *testing.T) {
	resolvers := &Stub{}

	srv := handler.New(NewExecutableSchema(Config{Resolvers: resolvers}))
	srv.AddTransport(transport.POST{})
	c := client.New(srv)

	resolvers.QueryResolver.Fallback = func(ctx context.Context, arg FallbackToStringEncoding) (FallbackToStringEncoding, error) {
		return arg, nil
	}

	t.Run("fallback to string passthrough", func(t *testing.T) {
		var resp struct {
			Fallback string
		}
		c.MustPost(`query { fallback(arg: A) }`, &resp)
		require.Equal(t, "A", resp.Fallback)
	})
}
