package subdir

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ssmccoy/gqlgen/_examples/embedding/subdir/gendir"
	"github.com/ssmccoy/gqlgen/client"
	"github.com/ssmccoy/gqlgen/graphql/handler"
	"github.com/ssmccoy/gqlgen/graphql/handler/transport"
)

func TestEmbeddingWorks(t *testing.T) {
	srv := handler.New(NewExecutableSchema(Config{Resolvers: &Resolver{}}))
	srv.AddTransport(transport.POST{})
	c := client.New(srv)
	var resp struct {
		InSchemadir string
		Parentdir   string
		Subdir      string
	}
	c.MustPost(`{
				inSchemadir
				parentdir
				subdir
			}
		`, &resp)

	require.Equal(t, "example", resp.InSchemadir)
	require.Equal(t, "example", resp.Parentdir)
	require.Equal(t, "example", resp.Subdir)
}

func TestEmbeddingWorksInGendir(t *testing.T) {
	srv := handler.New(gendir.NewExecutableSchema(gendir.Config{Resolvers: &GendirResolver{}}))
	srv.AddTransport(transport.POST{})
	c := client.New(srv)
	var resp struct {
		InSchemadir string
		Parentdir   string
		Subdir      string
	}
	c.MustPost(`{
				inSchemadir
				parentdir
				subdir
			}
		`, &resp)

	require.Equal(t, "example", resp.InSchemadir)
	require.Equal(t, "example", resp.Parentdir)
	require.Equal(t, "example", resp.Subdir)
}
