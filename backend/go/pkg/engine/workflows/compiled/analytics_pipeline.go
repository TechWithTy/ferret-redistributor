// Auto-generated workflow: Analytics Pipeline
package workflows

import (
    "fmt"
    "github.com/bitesinbyte/ferret/pkg/engine/auth"
    "github.com/bitesinbyte/ferret/pkg/engine/cache"
    "github.com/bitesinbyte/ferret/pkg/engine/factory"
    "github.com/bitesinbyte/ferret/pkg/engine/telemetry"
)

func RunAnalyticsPipeline() {
	fmt.Println("Running Workflow: Analytics Pipeline")
	auth.JWTAuth{}.Authenticate()
	cache.RedisCache{}.Save("workflow", "active")

	factory.Node{Type: "trigger", Name: "Cron"}.Execute()
	factory.Node{Type: "action", Name: "Extract"}.Execute()
	factory.Node{Type: "action", Name: "Transform"}.Execute()
	factory.Node{Type: "action", Name: "Load"}.Execute()
	telemetry.Sentry{}.TrackEvent("workflow_completed")
	fmt.Println("Workflow complete.")
}
