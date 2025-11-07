FROM golang:1.25 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o flowboard-backend cmd/main.go

FROM gcr.io/distroless/base-debian11
WORKDIR /app
COPY --from=builder /app/flowboard-backend .
COPY .env .
EXPOSE 8080
CMD ["./flowboard-backend"]
