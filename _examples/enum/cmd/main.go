package main

import (
	"log"
	"net/http"

	"github.com/ssmccoy/gqlgen/_examples/enum/api"
	"github.com/ssmccoy/gqlgen/graphql/handler"
	"github.com/ssmccoy/gqlgen/graphql/handler/transport"
	"github.com/ssmccoy/gqlgen/graphql/playground"
)

func main() {
	srv := handler.New(
		api.NewExecutableSchema(api.Config{Resolvers: &api.Resolver{}}),
	)

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	http.Handle("/", playground.Handler("Enum", "/query"))
	http.Handle("/query", srv)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
