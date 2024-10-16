FROM golang:1.22.5 AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY .. .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o service ./cmd/app

FROM scratch

WORKDIR /root/
COPY --from=builder /app/service .

EXPOSE 8080
EXPOSE 50051

CMD ["./service"]
