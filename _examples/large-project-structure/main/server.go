package main

import (
	"log"
	"net/http"
	"os"

	"github.com/vektah/gqlparser/v2/ast"

	"github.com/ssmccoy/gqlgen/_examples/large-project-structure/integration"
	"github.com/ssmccoy/gqlgen/_examples/large-project-structure/main/graph"
	_ "github.com/ssmccoy/gqlgen/_examples/large-project-structure/shared"
	"github.com/ssmccoy/gqlgen/graphql/handler"
	"github.com/ssmccoy/gqlgen/graphql/handler/extension"
	"github.com/ssmccoy/gqlgen/graphql/handler/lru"
	"github.com/ssmccoy/gqlgen/graphql/handler/transport"
	"github.com/ssmccoy/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Create a new executable schema with the composed resolver
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			ExternalQueryResolver: &integration.Resolver{},
			// Add other team resolvers here
		},
	}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
