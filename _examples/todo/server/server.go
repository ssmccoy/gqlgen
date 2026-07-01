package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/ssmccoy/gqlgen/_examples/todo"
	"github.com/ssmccoy/gqlgen/graphql/handler"
	"github.com/ssmccoy/gqlgen/graphql/handler/transport"
	"github.com/ssmccoy/gqlgen/graphql/playground"
)

func main() {
	srv := handler.New(todo.NewExecutableSchema(todo.New()))
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetRecoverFunc(func(ctx context.Context, err any) (userMessage error) {
		// send this panic somewhere
		log.Print(err)
		debug.PrintStack()
		return errors.New("user message on panic")
	})

	http.Handle("/", playground.Handler("Todo", "/query"))
	http.Handle("/query", srv)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
