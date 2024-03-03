dev: styles templates api-run

build: clean styles templates-format templates api

run: build
	./main

api:
	go build ./cmd/main.go

api-run:
	go run ./cmd/main.go

styles:
	tailwindcss -i ./public/styles/styles.css -o ./public/assets/styles.css

styles-watch:
	tailwindcss -i ./public/styles/styles.css -o ./public/assets/styles.css --watch &

templates:
	templ generate -path=./public

templates-format:
	templ fmt ./public

templates-watch:
	templ generate -path=./public -watch

clean:
	rm -f main