FROM --platform=$TARGETPLATFORM registry.hub.docker.com/library/golang:1.22 as builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM

WORKDIR /build

RUN echo "hello world" && echo $TARGETPLATFORM && echo $BUILDPLATFORM
RUN go env
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build  -installsuffix cgo -ldflags="-s -w" -o /api .

# RUN useradd -M --uid 1000 radix-log-api

FROM --platform=$TARGETPLATFORM docker.io/library/alpine

# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /api /usr/local/bin/api
# COPY --from=builder /etc/passwd /etc/passwd
USER 1000

ENTRYPOINT ["/usr/local/bin/api"]
