# Database Migrations in AWS

This document describes how database migrations are managed and executed in AWS environments.

## Strategy

We use a "One-Off Task" strategy for migrations. Instead of running migrations automatically on service startup (which can be risky with multiple replicas), we execute them as a standalone ECS Task using the same container image as the API service.

### Components

1. **Migration Binary**: A Go CLI tool located in `cmd/migrate` that uses `golang-migrate`.
2. **Container Image**: The standard API Docker image includes the `migrate` binary and the `migrations/` SQL files.
3. **ECS Task Definition**: The same task definition used by the API service, but we override the `command` when running it for migrations.
4. **IAM Permissions**: The ECS Task Execution Role and Task Role must have permissions to access RDS (via IAM Auth if enabled) and pull the image from ECR.

## Terraform Suggestions

To support migrations in AWS, ensure your Terraform configuration includes the following.

### 1. IAM Policy for RDS IAM Authentication

If using `DB_IAM_AUTH=true`, the ECS Task Role needs permission to connect to the RDS instance.

```hcl
resource "aws_iam_policy" "rds_connect" {
  name        = "emmy-${var.env}-rds-connect"
  description = "Allow ECS task to connect to RDS via IAM"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action   = "rds-db:connect"
        Effect   = "Allow"
        Resource = "arn:aws:rds-db:${var.region}:${data.aws_caller_identity.current.account_id}:dbuser:${var.rds_resource_id}/${var.db_user}"
      },
    ]
  })
}

resource "aws_iam_role_policy_attachment" "api_rds_connect" {
  role       = aws_iam_role.api_task_role.name
  policy_arn = aws_iam_policy.rds_connect.arn
}
```

### 2. Security Group Rules

The ECS tasks must be able to reach the RDS instance on port 5432.

```hcl
resource "aws_security_group_rule" "ecs_to_rds" {
  type                     = "ingress"
  from_port               = 5432
  to_port                 = 5432
  protocol                = "tcp"
  security_group_id       = aws_security_group.rds.id
  source_security_group_id = aws_security_group.ecs_api.id
}
```

## Running Migrations

### Locally (Dev)

```bash
make migrate-up
```

### Remote (AWS)

To run migrations in a specific environment (e.g., `dev`, `test`, `prod`):

```bash
make migrate-remote ENV=dev
```

This command uses the `scripts/migrate-rds` script, which:

1. Finds the ECS cluster and service for the environment.
2. Identifies the subnets and security groups used by the running service.
3. Starts a new Fargate task with a command override: `./migrate -path ./migrations up`.
4. Waits for the task to complete and reports the exit code.

## Handling Failures

If a migration fails and leaves the database in a "dirty" state, you may need to manually fix the database and then use the `force` command to reset the version.

```bash
# Example: force version to 1
make migrate-remote ENV=dev COMMAND="force 1"
```
