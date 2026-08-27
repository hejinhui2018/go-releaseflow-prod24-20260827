FROM golang:1.23.12-alpine
WORKDIR /src
COPY . .
RUN go test ./... && go vet ./... && go build ./...
ENTRYPOINT ["/bin/sh"]

