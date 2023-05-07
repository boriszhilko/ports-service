FROM golangci/golangci-lint:v1.42.1-alpine AS lint
WORKDIR /go/src/github.com/boriszhilko/ports-service
COPY . .
RUN golangci-lint run ./...

FROM golang:1.18 AS build
WORKDIR /go/src/github.com/boriszhilko/ports-service
COPY . .

RUN go mod download
RUN go build -o ./out/ports-service ./cmd

# Production image
FROM alpine:3.14

RUN apk --no-cache add libc6-compat
COPY --from=build /go/src/github.com/boriszhilko/ports-service/out/ports-service /usr/local/bin/
COPY assets/ports.json assets/ports.json

EXPOSE 8080

CMD ["ports-service"]