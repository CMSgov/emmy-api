FROM --platform=$BUILDPLATFORM artifactory.cloud.cms.gov/docker-remote/library/golang:1.25-alpine AS builder

ARG TARGETOS
ARG TARGETARCH
ARG ORCHESTRION_VERSION=v1.9.0

WORKDIR /build

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download && GOSUMDB=off go install github.com/DataDog/orchestrion@"${ORCHESTRION_VERSION}"

COPY . .

ENV CGO_ENABLED=0
RUN GOSUMDB=off GOOS="${TARGETOS:-linux}" GOARCH="$TARGETARCH" orchestrion go build -ldflags="-s -w" -a -o apiserver . && \
    GOOS="${TARGETOS:-linux}" GOARCH="$TARGETARCH" go build -ldflags="-s -w" -a -o migrate ./cmd/migrate

FROM artifactory.cloud.cms.gov/docker-remote/library/alpine:3.23

COPY --chmod=0755 --from=builder ["/build/apiserver", "/build/migrate", "/"]
COPY --chmod=0644 migrations/ /migrations/

RUN apk add --no-cache dumb-init ca-certificates

EXPOSE 8000

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/apiserver"]
