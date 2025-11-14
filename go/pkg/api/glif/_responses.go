package glif

import "encoding/json"

// RunWorkflowResponse captures the Simple API response.
type RunWorkflowResponse struct {
	ID         string                     `json:"id"`
	Inputs     map[string]json.RawMessage `json:"inputs"`
	Output     string                     `json:"output"`
	OutputFull json.RawMessage            `json:"outputFull"`
	Price      string                     `json:"price"`
	Nodes      []json.RawMessage          `json:"nodes"`
	Error      string                     `json:"error"`
}

// Workflow represents the metadata returned from the /api/glifs endpoints.
type Workflow struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Output      string       `json:"output"`
	OutputType  string       `json:"outputType"`
	User        WorkflowUser `json:"user"`
	Data        WorkflowData `json:"data"`
}

// WorkflowUser contains the workflow author's public profile fields.
type WorkflowUser struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Image    string `json:"image"`
}

// WorkflowData exposes the workflow nodes returned by the API.
type WorkflowData struct {
	Nodes []WorkflowNode `json:"nodes"`
}

// WorkflowNode describes an individual node/block configured inside the workflow graph.
type WorkflowNode struct {
	Name   string                     `json:"name"`
	Type   string                     `json:"type"`
	Params map[string]json.RawMessage `json:"params"`
}

