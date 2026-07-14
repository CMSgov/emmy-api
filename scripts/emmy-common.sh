#!/usr/bin/env bash
set -euo pipefail

readonly EMMY_IMAGE_REPO="${EMMY_IMAGE_REPO:-}"
readonly EMMY_VALID_ENVS=("dev" "test" "demo" "uat" "sandbox" "prod")

emmy_repo_root() {
  cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd
}

emmy_require_cmd() {
  local cmd="$1"

  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "error: required command not found: $cmd" >&2
    exit 1
  fi
}

emmy_validate_env() {
  local env="$1"
  local valid_env

  for valid_env in "${EMMY_VALID_ENVS[@]}"; do
    if [[ "$env" == "$valid_env" ]]; then
      return 0
    fi
  done

  echo "error: invalid ENV '$env'; expected one of: ${EMMY_VALID_ENVS[*]}" >&2
  exit 1
}

emmy_image_tag() {
  if [[ -n "${IMAGE_TAG:-}" ]]; then
    printf '%s\n' "$IMAGE_TAG"
    return 0
  fi

  git -C "$(emmy_repo_root)" rev-parse --short HEAD
}

emmy_image_uri() {
  local tag="${1:-$(emmy_image_tag)}"

  if [[ -z "$EMMY_IMAGE_REPO" ]]; then
    echo "error: EMMY_IMAGE_REPO must be set to the ECR repository URI" >&2
    exit 1
  fi

  printf '%s:%s\n' "$EMMY_IMAGE_REPO" "$tag"
}

emmy_task_family() {
  local env="$1"
  printf 'emmy-%s-api\n' "$env"
}

emmy_service_name() {
  local env="$1"
  printf 'emmy-%s-api\n' "$env"
}

emmy_cluster_name() {
  local env="$1"
  printf 'emmy-%s\n' "$env"
}
