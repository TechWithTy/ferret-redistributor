// Auto-generated workflow: Lead Import
package workflows

import (
    "fmt"
    "github.com/bitesinbyte/ferret/pkg/engine/auth"
    "github.com/bitesinbyte/ferret/pkg/engine/cache"
    "github.com/bitesinbyte/ferret/pkg/engine/factory"
    "github.com/bitesinbyte/ferret/pkg/engine/queue"
    "github.com/bitesinbyte/ferret/pkg/engine/telemetry"
)

func RunLeadImport() {
	fmt.Println("Running Workflow: Lead Import")
	auth.JWTAuth{}.Authenticate()
	cache.RedisCache{}.Save("workflow", "active")

	factory.Node{Type: "trigger", Name: "Webhook"}.Execute()
	factory.Node{Type: "action", Name: "Transform Lead"}.Execute()
	factory.Node{Type: "action", Name: "Push to CRM"}.Execute()
	_ = queue.NatsEngine{}
	factory.Node{Type: "on_error", Name: "Notify"}.Execute()
	telemetry.Sentry{}.TrackEvent("workflow_completed")
	fmt.Println("Workflow complete.")
}
