package main

import (
	"log"
	"net/http"

	unionextension "github.com/ssmccoy/gqlgen/_examples/union-extension"
	"github.com/ssmccoy/gqlgen/graphql/handler"
	"github.com/ssmccoy/gqlgen/graphql/handler/transport"
	"github.com/ssmccoy/gqlgen/graphql/playground"
)

func main() {
	srv := handler.New(
		unionextension.NewExecutableSchema(
			unionextension.Config{Resolvers: &unionextension.Resolver{}},
		),
	)
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	http.Handle("/", playground.Handler("Union Extension Demo", "/query"))
	http.Handle("/query", srv)
	log.Fatal(http.ListenAndServe(":8086", nil))
}
