FROM golang:1.14 as builder

ARG VERSION
ARG COMMIT

ENV GO111MODULE=on \
    GOPROXY=https://proxy.golang.org

WORKDIR /go/src/github.com/summerwind/h2spec
COPY go.mod go.sum .
RUN go mod download

COPY . /workspace
WORKDIR /workspace

RUN go vet ./...
RUN go test -v ./...
RUN CGO_ENABLED=0 go build -ldflags "-X main.VERSION=${VERSION} -X main.COMMIT=${COMMIT}" ./cmd/h2spec

###################

FROM ubuntu:18.04

COPY --from=builder /workspace/h2spec /usr/local/bin/h2spec

ENTRYPOINT ["/usr/local/bin/h2spec"]
