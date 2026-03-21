#!/bin/bash

ECS_CLUSTER="emmy-dev"
ECS_SERVICE="emmy-dev-api"
TASK_DEFINITION_FAMILY="emmy-dev-api"
NEW_IMAGE_URI="artifactory.cloud.cms.gov/emmy-docker/emmy-api:latest"

# 1. Fetch the existing task definition as a template
echo "Fetching current task definition for family: ${TASK_DEFINITION_FAMILY}..."
TASK_DEF_JSON=$(aws ecs describe-task-definition --task-definition "${TASK_DEFINITION_FAMILY}" | jq '.taskDefinition')

# 2. Modify the JSON: update the image and remove non-creatable keys
# Task definitions are immutable; you must create a new revision.
echo "Updating image URI to: ${NEW_IMAGE_URI} and cleaning up keys..."
UPDATED_TASK_DEF_JSON=$(echo "${TASK_DEF_JSON}" | \
                            jq '.containerDefinitions[0].image = "'"${NEW_IMAGE_URI}"'"' | \
                            jq 'del(.compatibilities, .status, .taskDefinitionArn, .revision, .requiresAttributes, .registeredAt, .registeredBy)')

# 3. Register the new task definition revision
echo "Registering new task definition revision..."
NEW_REVISION_ARN=$(aws ecs register-task-definition --cli-input-json "${UPDATED_TASK_DEF_JSON}" | jq -r '.taskDefinition.taskDefinitionArn')

if [ -z "$NEW_REVISION_ARN" ]; then
    echo "Failed to register new task definition."
    exit 1
fi

echo "New task definition revision registered: ${NEW_REVISION_ARN}"

# 4. Update the ECS service to use the new task definition
# The service scheduler will stop old tasks and start new ones with the new definition.
echo "Updating ECS service ${ECS_SERVICE} in cluster ${ECS_CLUSTER} to use the new revision..."
aws ecs update-service --cluster "${ECS_CLUSTER}" --service "${ECS_SERVICE}" --task-definition "${NEW_REVISION_ARN}" --force-new-deployment

echo "Service update initiated. ECS will now deploy new tasks."
# You can optionally add a 'wait for service stabilization' step here using 'aws ecs wait services-stable'
