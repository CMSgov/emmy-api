FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags="-s -w" -a -o apiserver .

FROM alpine:3.23

COPY --chmod=0755 --from=builder ["/build/apiserver", "/"]

RUN apk add --no-cache dumb-init ca-certificates

EXPOSE 8000

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/apiserver"]
