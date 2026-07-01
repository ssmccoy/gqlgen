package main

import (
	"log"
	"net/http"

	"github.com/ssmccoy/gqlgen/_examples/selection"
	"github.com/ssmccoy/gqlgen/graphql/handler"
	"github.com/ssmccoy/gqlgen/graphql/handler/transport"
	"github.com/ssmccoy/gqlgen/graphql/playground"
)

func main() {
	srv := handler.New(
		selection.NewExecutableSchema(selection.Config{Resolvers: &selection.Resolver{}}),
	)
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	http.Handle("/", playground.Handler("Selection Demo", "/query"))
	http.Handle("/query", srv)
	log.Fatal(http.ListenAndServe(":8086", nil))
}
