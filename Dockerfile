FROM golang:1.21.4-alpine as binary-builder
RUN apk add --no-cache build-base
WORKDIR /builder
RUN go install github.com/a-h/templ/cmd/templ@latest
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN templ generate
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build \
  -ldflags='-w -s' -a \
  -o app cmd/*.go

FROM alpine:3.18
ENV APP_PORT=3333
WORKDIR /app
COPY --from=binary-builder /builder/app .
EXPOSE $APP_PORT
ENTRYPOINT ["./app"]
