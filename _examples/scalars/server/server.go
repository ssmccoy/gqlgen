package main

import (
	"log"
	"net/http"

	"github.com/ssmccoy/gqlgen/_examples/scalars"
	"github.com/ssmccoy/gqlgen/graphql/handler"
	"github.com/ssmccoy/gqlgen/graphql/handler/transport"
	"github.com/ssmccoy/gqlgen/graphql/playground"
)

func main() {
	srv := handler.New(
		scalars.NewExecutableSchema(scalars.Config{Resolvers: &scalars.Resolver{}}),
	)
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	http.Handle("/", playground.Handler("Starwars", "/query"))
	http.Handle("/query", srv)

	log.Fatal(http.ListenAndServe(":8084", nil))
}
