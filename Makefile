SHELL := /usr/bin/env bash

.PHONY: help build push deploy release print-image-tag print-image-uri migrate-up migrate-down migrate-version migrate-remote

ENV ?= dev
CONTEXT ?= .
IMAGE_TAG ?=

help:
	@printf '%s\n' \
	  'Targets:' \
	  '  make build [CONTEXT=.] [PLATFORM=linux/amd64] Build the local image tagged with the short git SHA' \
	  '  make push [CONTEXT=.] [PLATFORMS=linux/amd64] Build and push the image to Artifactory' \
	  '  make deploy ENV=<env> [IMAGE_TAG=tag]  Update the ECS task definition and deploy it' \
	  '  make release ENV=<env> [CONTEXT=.]     Push and then deploy the current image tag' \
	  '  make print-image-tag                   Print the image tag that will be used' \
	  '  make print-image-uri                   Print the fully-qualified image URI' \
	  '  make migrate-up                        Apply all pending migrations' \
	  '  make migrate-down                      Revert the last migration' \
	  '  make migrate-version                   Show the current migration version' \
	  '  make migrate-remote ENV=<env> [COMMAND=up] Run migrations on an AWS environment' \
	  '' \
	  'Valid ENV values: dev test demo uat sandbox prod'

build:
	@IMAGE_TAG='$(IMAGE_TAG)' ./scripts/build-image "$(CONTEXT)"

push:
	@IMAGE_TAG='$(IMAGE_TAG)' ./scripts/push-image "$(CONTEXT)"

deploy:
	@if [[ -z "$(ENV)" ]]; then \
	  echo "error: ENV is required; valid values: dev test demo uat sandbox prod" >&2; \
	  exit 1; \
	fi
	@IMAGE_TAG='$(IMAGE_TAG)' ./scripts/deploy-ecs "$(ENV)"

release:
	@if [[ -z "$(ENV)" ]]; then \
	  echo "error: ENV is required; valid values: dev test demo uat sandbox prod" >&2; \
	  exit 1; \
	fi
	@IMAGE_TAG='$(IMAGE_TAG)' ./scripts/push-image "$(CONTEXT)"
	@IMAGE_TAG='$(IMAGE_TAG)' ./scripts/deploy-ecs "$(ENV)"

print-image-tag:
	@IMAGE_TAG='$(IMAGE_TAG)' bash -lc 'source ./scripts/emmy-common.sh && emmy_image_tag'

print-image-uri:
	@IMAGE_TAG='$(IMAGE_TAG)' bash -lc 'source ./scripts/emmy-common.sh && emmy_image_uri'

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down

migrate-version:
	@go run cmd/migrate/main.go version

migrate-remote:
	@if [[ -z "$(ENV)" ]]; then \
	  echo "error: ENV is required; valid values: dev test demo uat sandbox prod" >&2; \
	  exit 1; \
	fi
	@./scripts/migrate-rds "$(ENV)" "$(COMMAND)"
