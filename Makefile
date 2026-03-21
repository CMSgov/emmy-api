TAG ?= latest
IMAGE ?= artifactory.cloud.cms.gov/emmy-docker/emmy-api

TASK_DEFINITION = arn:aws:ecs:us-east-1:554269192106:task-definition/emmy-dev-api
SERVICE_NAME = emmy-dev-api
CLUSTER_NAME = emmy-dev

.PHONY: build push tag update-task-def deploy all

all: deploy

tag:
	echo $(TAG)

build:
	docker buildx build -t $(IMAGE):$(TAG) --platform linux/amd64 --push .

deploy: build
	./scripts/deploy-task.sh
