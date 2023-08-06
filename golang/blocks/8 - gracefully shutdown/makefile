VERSION := 1.0.0

tidy:
	go mod tidy && go mod vendor

run:
	go run main.go

build:
	go build -ldflags="-X main.version=$(VERSION)" main.go