# Multi-state build. 
# Build stage
FROM golang:rc-alpine
RUN apk add make git

WORKDIR /go/src/github.com/roaanv/k8sti

COPY . .
RUN make

# REAL image
FROM alpine:latest

COPY --from=0 /go/src/github.com/roaanv/k8sti/main /usr/local/k8sti

ENTRYPOINT ["/usr/local/k8sti"]
