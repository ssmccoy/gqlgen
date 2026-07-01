bin_name=server

build:
	@echo "building binary..."
# go generate gives missing go sum entry for module errors
# https://github.com/ssmccoy/gqlgen/issues/1483
# you will need to first do a go get -u github.com/ssmccoy/gqlgen
	go run -mod=mod github.com/ssmccoy/gqlgen generate .
	go build -o ${bin_name} server.go
