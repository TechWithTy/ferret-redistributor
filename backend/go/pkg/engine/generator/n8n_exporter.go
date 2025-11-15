package generator

import (
    "encoding/json"
    "fmt"
    "strings"
)

// ExportToN8N converts pseudo workflow to a minimal n8n JSON.
func ExportToN8N(pseudoPath, outPath string) error {
    wf := ParsePseudo(pseudoPath)
    var nodes []map[string]any

    for i, n := range wf.NodeMatches {
        nType, nName := n[1], n[2]
        node := map[string]any{
            "id":   fmt.Sprintf("%s_%d", strings.ReplaceAll(nName, " ", "_"), i+1),
            "name": nName,
            "type": map[string]string{
                "trigger":  "n8n-nodes-base.webhook",
                "action":   "n8n-nodes-base.httpRequest",
                "switch":   "n8n-nodes-base.switch",
                "merge":    "n8n-nodes-base.merge",
                "subflow":  "n8n-nodes-base.executeWorkflow",
                "on_error": "n8n-nodes-base.code",
            }[nType],
        }
        nodes = append(nodes, node)
    }

    // Build simple connections map: source name -> targets
    connections := map[string]any{}
    for _, c := range wf.Connections {
        from := c[1]
        targetsRaw := c[2]
        var targets []string
        if strings.HasPrefix(targetsRaw, "[") {
            // strip [ ] and split
            inner := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(targetsRaw, "["), "]"))
            if inner != "" {
                parts := strings.Split(inner, ",")
                for _, p := range parts {
                    t := strings.Trim(strings.TrimSpace(p), "\"")
                    if t != "" {
                        targets = append(targets, t)
                    }
                }
            }
        } else {
            targets = []string{strings.Trim(targetsRaw, "\"")}
        }
        // n8n format: { "main": [[ {"node": "Target"}, ... ]] }
        var connRow []map[string]any
        for _, t := range targets {
            connRow = append(connRow, map[string]any{"node": t})
        }
        connections[from] = map[string]any{"main": [][]map[string]any{connRow}}
    }

    out := map[string]any{
        "name":        wf.Name,
        "nodes":       nodes,
        "connections": connections,
    }

    b, err := json.MarshalIndent(out, "", "  ")
    if err != nil {
        return err
    }
    return writeFile(outPath, b)
}
