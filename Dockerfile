# stage build
FROM golang:1.23-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /tg_echo ./cmd/tg_echo

# stage runtime
FROM alpine:3.20
ENV PORT=8080
COPY --from=build /tg_echo /tg_echo
RUN chmod +x /tg_echo
ENTRYPOINT ["/tg_echo"]