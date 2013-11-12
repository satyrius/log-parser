build:
	go build -o parser parser.go

build.linux:
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o parser.linux parser.go
