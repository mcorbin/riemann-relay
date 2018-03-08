
testv:
	go test -v ./...

test:
	go test ./...

build:
	go build github.com/mcorbin/riemann-relay/cmd/relay
