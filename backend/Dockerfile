FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/api-server/main.go

FROM gcr.io/distroless/static-debian12
COPY --from=builder /server /server
COPY --from=builder /app/config.yaml /config.yaml
WORKDIR /
EXPOSE 8080
ENTRYPOINT ["/server"]
