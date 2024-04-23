FROM registry.hub.docker.com/library/golang:1.22 as builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build  -installsuffix cgo -ldflags="-s -w" -o /api .

# RUN useradd -M --uid 1000 radix-log-api

FROM alpine
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /api /usr/local/bin/api
# COPY --from=builder /etc/passwd /etc/passwd
USER 1000

ENTRYPOINT ["/usr/local/bin/api"]
