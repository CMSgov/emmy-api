FROM --platform=$BUILDPLATFORM artifactory.cloud.cms.gov/docker-remote/library/golang:1.25-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download && go install github.com/DataDog/orchestrion@latest

COPY . .

ENV CGO_ENABLED=0
RUN GOOS="${TARGETOS:-linux}" GOARCH="$TARGETARCH" orchestrion go build -ldflags="-s -w" -a -o apiserver . && \
    GOOS="${TARGETOS:-linux}" GOARCH="$TARGETARCH" go build -ldflags="-s -w" -a -o migrate ./cmd/migrate

FROM artifactory.cloud.cms.gov/docker-remote/library/alpine:3.23

COPY --chmod=0755 --from=builder ["/build/apiserver", "/build/migrate", "/"]
COPY --chmod=0644 migrations/ /migrations/

RUN apk add --no-cache dumb-init ca-certificates

EXPOSE 8000

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/apiserver"]
