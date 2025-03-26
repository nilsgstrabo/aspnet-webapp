FROM golang:1.24-alpine3.21x AS builder
# FROM golang:1.22-alpine3.20 AS builder
# FROM docker.io/nilsgustavstrabo/golang:1.22-alpine3.20 as builder
# FROM --platform=$TARGETPLATFORM docker.io/library/golang:1.22-alpine3.20 as builder


# ARG TARGETPLATFORM
# ARG BUILDPLATFORM

WORKDIR /build

# RUN --mount=type=secret,id=SECRET1,dst=/abc/my-secrets/secret-1.txt export BUILD_ARG=$(cat /abc/my-secrets/secret-1.txt) && \
#     echo $BUILD_ARG

RUN --mount=type=secret,id=SECRET1,dst=/abc/my-secrets/secret-1.txt --mount=type=secret,id=SECRET2a,dst=/abc/my-secrets/secret-2.txt \
    export BUILD_ARG=$(cat /abc/my-secrets/secret-1.txt) && \
    export BUILD_ARG2=$(cat /abc/my-secrets/secret-2.txt) && \
    echo $BUILD_ARG and $BUILD_ARG2 && \
    echo $BUILD_ARG and $BUILD_ARG2 > /mysecrets.txt

# RUN echo "hello world" && echo $TARGETPLATFORM && echo $BUILDPLATFORM
RUN echo "hello world"
RUN go env
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build  -installsuffix cgo -ldflags="-s -w" -o /api .

# RUN useradd -M --uid 1000 radix-log-api

# FROM --platform=$TARGETPLATFORM docker.io/library/alpine:3.20
FROM alpine:3.21

# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /api /usr/local/bin/api
COPY --from=builder /mysecrets.txt /mysecrets.txt
# COPY --from=builder /etc/passwd /etc/passwd
USER 3000

ENTRYPOINT ["/usr/local/bin/api"]
