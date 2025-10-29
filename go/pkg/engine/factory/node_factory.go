package factory

import "fmt"

type Node struct {
    Type       string
    Name       string
    Parameters map[string]string
}

func (n Node) Execute() {
    fmt.Printf("[%s] Node executed: %s\n", n.Type, n.Name)
}

