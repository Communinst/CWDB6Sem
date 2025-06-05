FROM golang:1.23.2-alpine as builder
WORKDIR /build
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /main cmd/main.go

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /main /main
#COPY config ./config
#ENV CONFIG_PATH=/config/config.yaml
ENTRYPOINT ["/main"]