package factory

type Workflow struct {
    Name  string
    Nodes []Node
}

func NewWorkflow(name string, nodes []Node) Workflow {
    return Workflow{Name: name, Nodes: nodes}
}

