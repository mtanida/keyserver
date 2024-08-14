#----------------------------------------------------------------------
# Builder image
FROM golang:1.23.0-alpine3.20 AS builder

WORKDIR /app
COPY go.mod go.sum main.go  ./
RUN go mod download
RUN CGO_ENABLED=0 go build -o keyserver main.go

#----------------------------------------------------------------------
# Final Runtime image
FROM alpine:3.20 AS runtime
WORKDIR /app
COPY --from=builder /app/keyserver keyserver
ENTRYPOINT ["./keyserver"]
