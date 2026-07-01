package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ssmccoy/gqlgen/graphql/handler"
	"github.com/ssmccoy/gqlgen/graphql/handler/debug"
	"github.com/ssmccoy/gqlgen/graphql/handler/transport"
	"github.com/ssmccoy/gqlgen/graphql/playground"
	"github.com/ssmccoy/gqlgen/plugin/federation/testdata/computedrequires"
	"github.com/ssmccoy/gqlgen/plugin/federation/testdata/computedrequires/generated"
)

const defaultPort = "4003"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.New(
		generated.NewExecutableSchema(generated.Config{Resolvers: &computedrequires.Resolver{}}),
	)
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.Use(&debug.Tracer{})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
