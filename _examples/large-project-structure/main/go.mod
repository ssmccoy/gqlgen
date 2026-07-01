module github.com/ssmccoy/gqlgen/_examples/large-project-structure/main

go 1.25.0

require (
	github.com/ssmccoy/gqlgen v0.17.89
	github.com/ssmccoy/gqlgen/_examples/large-project-structure/integration v0.0.0-00010101000000-000000000000
	github.com/ssmccoy/gqlgen/_examples/large-project-structure/shared v0.0.0
	github.com/vektah/gqlparser/v2 v2.5.34
)

replace github.com/ssmccoy/gqlgen/_examples/large-project-structure/shared => ../shared

replace github.com/ssmccoy/gqlgen/_examples/large-project-structure/integration => ../integration

require (
	github.com/agnivade/levenshtein v1.2.1 // indirect
	github.com/coder/websocket v1.8.14 // indirect
	github.com/go-viper/mapstructure/v2 v2.5.0 // indirect
	github.com/goccy/go-yaml v1.19.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/sosodev/duration v1.4.0 // indirect
	github.com/urfave/cli/v3 v3.9.0 // indirect
	golang.org/x/mod v0.36.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/text v0.37.0 // indirect
	golang.org/x/tools v0.45.0 // indirect
)

replace github.com/ssmccoy/gqlgen => ../../../
