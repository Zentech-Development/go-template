build:
	go build -o ./APPNAME ./main.go

dev:
	air

clean:
	rm -f ./APPNAME
	rm -rf ./tmp