FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o service main.go

FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/service .

EXPOSE 8000
ENTRYPOINT ["./service"]
