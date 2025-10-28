package workers

import (
    "context"
    "strings"

    ig "github.com/bitesinbyte/ferret/pkg/external/instagram"
    "github.com/bitesinbyte/ferret/pkg/generator"
)

type TriggerMatcher struct {
    Triggers []string
}

func (m TriggerMatcher) Match(text string) (matched string, ok bool) {
    t := strings.ToLower(text)
    for _, trg := range m.Triggers {
        trg = strings.TrimSpace(strings.ToLower(trg))
        if trg == "" { continue }
        if strings.Contains(t, trg) {
            return trg, true
        }
    }
    return "", false
}

type DMWorker struct {
    IG        *ig.Client
    Matcher   TriggerMatcher
    Generator generator.AIMLGenerator
}

// ProcessTriggerComments scans comments on a media, detects trigger words,
// and attempts to DM the commenter with generated content.
func (w DMWorker) ProcessTriggerComments(ctx context.Context, mediaID string, prompt string) error {
    comments, err := w.IG.ListComments(ctx, mediaID)
    if err != nil { return err }
    for _, c := range comments {
        // Ignore self comments
        if c.FromID == w.IG.IGUserID() { continue }
        if _, ok := w.Matcher.Match(c.Text); !ok { continue }
        dmText, err := w.Generator.GenerateDM(ctx, prompt, map[string]string{
            "comment_text": c.Text,
            "username":     c.Username,
            "user_id":      c.FromID,
            "media_id":     mediaID,
        })
        if err != nil { return err }
        // SendDM may be unsupported unless messaging is configured
        if err := w.IG.SendDM(ctx, c.FromID, dmText); err != nil {
            return err
        }
    }
    return nil
}
