# Build the manager binary
FROM golang:1.22 AS builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} go build -a -o gameserver

FROM alpine:latest
WORKDIR /
RUN apk --no-cache add bash curl
COPY --from=builder /workspace/gameserver .
COPY check-idle.sh /check-idle.sh
COPY set-idle.sh /set-idle.sh
COPY set-busy.sh /set-busy.sh
USER 65532:65532

ENTRYPOINT ["/gameserver"]

