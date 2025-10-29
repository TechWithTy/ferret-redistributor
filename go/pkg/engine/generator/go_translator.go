package generator

import (
    "fmt"
    "strings"
)

// CompileToGo generates a simple Go program from the pseudo workflow.
func CompileToGo(pseudoPath, outPath string) error {
    wf := ParsePseudo(pseudoPath)
    var code strings.Builder

    code.WriteString(fmt.Sprintf("// Auto-generated workflow: %s\n", wf.Name))
    code.WriteString("package main\n\n")
    code.WriteString("import (\n")
    code.WriteString("\t\"fmt\"\n")
    code.WriteString("\t\"github.com/bitesinbyte/ferret/pkg/engine/auth\"\n")
    code.WriteString("\t\"github.com/bitesinbyte/ferret/pkg/engine/cache\"\n")
    code.WriteString("\t\"github.com/bitesinbyte/ferret/pkg/engine/factory\"\n")
    code.WriteString("\t\"github.com/bitesinbyte/ferret/pkg/engine/queue\"\n")
    code.WriteString("\t\"github.com/bitesinbyte/ferret/pkg/engine/telemetry\"\n")
    code.WriteString("\t\"github.com/bitesinbyte/ferret/pkg/engine/workers\"\n")
    code.WriteString(")\n\n")

    code.WriteString("func main() {\n")
    code.WriteString(fmt.Sprintf("\tfmt.Println(\"Running Workflow: %s\")\n", wf.Name))
    code.WriteString("\tauth.JWTAuth{}.Authenticate()\n")
    code.WriteString("\tcache.RedisCache{}.Save(\"workflow\", \"active\")\n\n")

    for _, node := range wf.NodeMatches {
        nodeType := node[1]
        nodeName := node[2]
        code.WriteString(fmt.Sprintf("\tfactory.Node{Type: \"%s\", Name: \"%s\"}.Execute()\n", nodeType, nodeName))
        // small showcase calls so imports are used
        if nodeType == "action" && strings.Contains(strings.ToLower(nodeName), "crm") {
            code.WriteString("\t_ = queue.NatsEngine{}\n")
        }
        if nodeType == "action" && strings.Contains(strings.ToLower(nodeName), "enrich") {
            code.WriteString("\t_ = workers.AIWorker{}\n")
        }
    }

    code.WriteString("\ttelemetry.Sentry{}.TrackEvent(\"workflow_completed\")\n")
    code.WriteString("\tfmt.Println(\"Workflow complete.\")\n")
    code.WriteString("}\n")

    return writeFile(outPath, []byte(code.String()))
}

