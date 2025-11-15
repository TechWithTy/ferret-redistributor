// Auto-generated workflow: CRM Sync
package workflows

import (
	"fmt"

	"github.com/bitesinbyte/ferret/pkg/engine/auth"
	"github.com/bitesinbyte/ferret/pkg/engine/cache"
	"github.com/bitesinbyte/ferret/pkg/engine/factory"
	"github.com/bitesinbyte/ferret/pkg/engine/queue"
	"github.com/bitesinbyte/ferret/pkg/engine/telemetry"
)

func RunCRMSync() {
	fmt.Println("Running Workflow: CRM Sync")
	auth.JWTAuth{}.Authenticate()
	cache.RedisCache{}.Save("workflow", "active")

	factory.Node{Type: "trigger", Name: "Schedule"}.Execute()
	factory.Node{Type: "action", Name: "Pull Updates"}.Execute()
	factory.Node{Type: "action", Name: "Push to CRM"}.Execute()
	_ = queue.NatsEngine{}
	telemetry.Sentry{}.TrackEvent("workflow_completed")
	fmt.Println("Workflow complete.")
}
