build:
	go build -o ./APPNAME ./cmd/server/main.go

run:
	./APPNAME

dev:
	air

clean:
	rm -f ./APPNAME
	rm -rf ./tmp