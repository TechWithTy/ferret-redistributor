package generator

import (
    "bufio"
    "os"
    "regexp"
    "strings"
)

type Workflow struct {
    Name        string
    Version     string
    NodeMatches [][]string // [full, type, name, using, path]
    Connections [][]string // [full, from, targetsRaw]
}

// ParsePseudo parses a minimal EOWR pseudo workflow file.
func ParsePseudo(path string) Workflow {
    data, _ := os.ReadFile(path)
    content := string(data)

    nameRe := regexp.MustCompile(`(?m)^\s*workflow\s+"([^"]+)"\s+version\s+([^:]+):`)
    nameMatch := nameRe.FindStringSubmatch(content)
    wf := Workflow{}
    if len(nameMatch) >= 3 {
        wf.Name = strings.TrimSpace(nameMatch[1])
        wf.Version = strings.TrimSpace(nameMatch[2])
    }

    nodeLineRe := regexp.MustCompile(`^\s*(trigger|action|switch|merge|subflow|on_error)\s+"([^"]+)"(?:\s+using\s+"([^"]+)")?(?:\s+at\s+"([^"]+)")?`)

    var nodes [][]string
    scanner := bufio.NewScanner(strings.NewReader(content))
    for scanner.Scan() {
        line := scanner.Text()
        if m := nodeLineRe.FindStringSubmatch(line); len(m) > 0 {
            using, pathAttr := "", ""
            if len(m) >= 4 {
                using = m[3]
            }
            if len(m) >= 5 {
                pathAttr = m[4]
            }
            nodes = append(nodes, []string{m[0], m[1], m[2], using, pathAttr})
        }
    }
    wf.NodeMatches = nodes

    connRe := regexp.MustCompile(`(?m)^\s*connect\s+"([^"]+)"\s*->\s*(\[[^\]]+\]|"[^"]+")`)
    wf.Connections = connRe.FindAllStringSubmatch(content, -1)

    return wf
}

