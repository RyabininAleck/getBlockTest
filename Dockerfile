FROM golang:1.22.4 AS build

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache libc6-compat

COPY --from=build /app/main .
COPY --from=build /app/config.json .

EXPOSE 8080

CMD ["./main"]
