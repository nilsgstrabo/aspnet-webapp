FROM docker.io/library/golang:1.22-alpine3.20 as builder
# FROM docker.io/nilsgustavstrabo/golang:1.22-alpine3.20 as builder
# FROM --platform=$TARGETPLATFORM docker.io/library/golang:1.22-alpine3.20 as builder


# ARG TARGETPLATFORM
# ARG BUILDPLATFORM

WORKDIR /build

RUN --mount=type=secret,id=SECRET1,target=/abc/my-secrets/secret-1.txt export BUILD_ARG=$(cat /abc/my-secrets/secret-1.txt) && \
    echo $BUILD_ARG

# RUN echo "hello world" && echo $TARGETPLATFORM && echo $BUILDPLATFORM
RUN echo "hello world"
RUN go env
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build  -installsuffix cgo -ldflags="-s -w" -o /api .

# RUN useradd -M --uid 1000 radix-log-api

# FROM --platform=$TARGETPLATFORM docker.io/library/alpine:3.20
FROM docker.io/library/alpine:3.20

# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /api /usr/local/bin/api
# COPY --from=builder /etc/passwd /etc/passwd
USER 1000

ENTRYPOINT ["/usr/local/bin/api"]
