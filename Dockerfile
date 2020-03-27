FROM golang:1.13.3 AS builder

WORKDIR /csgo-starter

# Install dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the source code and build the binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/csgo-starter

FROM alpine:latest

RUN mkdir /app

# Copy the binary from previous build stage and run
COPY --from=builder /go/bin/csgo-starter /app
WORKDIR /app
ENTRYPOINT [ "./csgo-starter" ]
