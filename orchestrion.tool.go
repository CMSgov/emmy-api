// This file pins Orchestrion and Datadog tracer integrations for repeatable builds.
//go:build tools

package tools

import (
	_ "github.com/DataDog/dd-trace-go/contrib/log/slog/v2"
	_ "github.com/DataDog/dd-trace-go/contrib/redis/go-redis.v9/v2"
	_ "github.com/DataDog/orchestrion"
)
