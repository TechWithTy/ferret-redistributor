package glif

// Visibility controls whether a workflow run is public or private.
type Visibility string

const (
	// VisibilityPrivate is the default for workflow runs.
	VisibilityPrivate Visibility = "PRIVATE"
	// VisibilityPublic allows Glif to expose the run publicly.
	VisibilityPublic Visibility = "PUBLIC"
)

// RunWorkflowRequest defines the payload sent to the Simple API.
type RunWorkflowRequest struct {
	// WorkflowID is required unless UsePathID is true and the ID is embedded in the URL.
	WorkflowID string
	// Inputs represents positional inputs in the order defined by the workflow.
	Inputs []string
	// NamedInputs maps internal block names to their values.
	NamedInputs map[string]string
	// Visibility overrides Glif's default of PRIVATE when provided.
	Visibility Visibility
	// Strict enforces strict mode (?strict=1) so missing inputs fail instead of falling back to defaults.
	Strict bool
	// UsePathID instructs the client to send the ID as part of the URL (https://simple-api.glif.app/<id>).
	UsePathID bool
}
