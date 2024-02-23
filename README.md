# Go Template

Base template for Go backend applications. To use, clone this repo, search for "go-template" and replace
with your application's name. 

## Installation
1. Install [Go](https://go.dev/doc/install)
1. Install [templ](https://templ.guide/quick-start/installation)
1. Install the [TailwindCSS](https://tailwindcss.com/blog/standalone-cli) standalone CLI

## Running 
1. Compile the templates
```sh
templ generate -path=./public
```
1. Compile the stylesheets
```sh
tailwindcss -i ./public/styles/styles.css -o ./public/assets/styles.css
```
1. Run the API
```sh
go run ./cmd/main.go
```

If you are working exclusively on the templates, include the watch option.
```sh
templ generate -path=./public -watch
tailwindcss -i ./public/styles/styles.css -o ./public/assets/styles.css --watch
```