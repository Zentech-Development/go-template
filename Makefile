build:
	go build -o ./APPNAME ./cmd/server/main.go

dev:
	air

clean:
	rm -f ./APPNAME
	rm -rf ./tmp