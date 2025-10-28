package generator

import "context"

// AIMLGenerator defines the minimal interface to produce DM text
// using the local ai-ml-models providers.
// Implementations can shell out to Python or call a local HTTP service.
type AIMLGenerator interface {
    GenerateDM(ctx context.Context, prompt string, variables map[string]string) (string, error)
}

// NoopGenerator is a fallback that returns a static message.
type NoopGenerator struct{}

func (NoopGenerator) GenerateDM(ctx context.Context, prompt string, variables map[string]string) (string, error) {
    _ = ctx
    _ = prompt
    _ = variables
    return "Thanks for your comment! I’ll DM you the details shortly.", nil
}

