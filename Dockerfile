# Build stage
ARG GOVERSION=1.25
ENV GOPROXY=direct
FROM golang:${GOVERSION} AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd .
RUN CGO_ENABLED=0 go build -o /autoheal .

# Runtime stage
FROM docker:cli

COPY --from=builder /autoheal /usr/local/bin/autoheal
USER root
ENTRYPOINT ["/usr/local/bin/autoheal"]
