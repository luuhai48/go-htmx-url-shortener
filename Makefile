dev:
	modd

run:
	@go run cmd/*.go

temp:
	@templ generate

tailwind:
	@tailwindcss -i ./style.pcss -o ./static/files/style.css --postcss

prepare: temp tailwind

build-app:
	CGO_ENABLED=1 GOARCH=amd64 go build -ldflags='-w -s' -a -o app cmd/*.go

build: prepare build-app	

kill:
	kill $(lsof -t -i tcp:3333)
