build:
	go build -o ./APPNAME ./cmd/server/main.go

run:
	./APPNAME -store=sqlite -binding=gin

dev:
	air

clean:
	rm -f ./APPNAME
	rm -rf ./tmp