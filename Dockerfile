FROM golang:1.24.2 AS builder
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download

# Build
COPY . .
RUN git rev-parse --short HEAD
RUN GIT_COMMIT=$(git rev-parse --short HEAD) && \
    CGO_ENABLED=0 go build -o hs -ldflags "-X main.GitCommit=${GIT_COMMIT}"

FROM alpine:latest
RUN apk update && apk add curl jq ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=builder /build/hs /bin/hs
EXPOSE 35444
CMD ["/bin/hs"]
