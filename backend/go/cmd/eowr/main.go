package main

import (
    "errors"
    "fmt"
    "os"
    "path/filepath"

    gen "github.com/bitesinbyte/ferret/pkg/engine/generator"
)

func usage() {
    fmt.Println("eowr - Engine-Oriented Workflow Runtime CLI")
    fmt.Println("Usage:")
    fmt.Println("  go run ./cmd/eowr compile <in.pseudo> <out.go>")
    fmt.Println("  go run ./cmd/eowr export  <in.pseudo> <out.json>")
    fmt.Println("  go run ./cmd/eowr compile-all")
    fmt.Println("  go run ./cmd/eowr export-all")
}

func compileOne(inPath, outPath string) error {
    if inPath == "" || outPath == "" {
        return errors.New("compile requires <in.pseudo> and <out.go>")
    }
    if err := gen.CompileToGo(inPath, outPath); err != nil {
        return err
    }
    fmt.Printf("Compiled: %s -> %s\n", inPath, outPath)
    return nil
}

func exportOne(inPath, outPath string) error {
    if inPath == "" || outPath == "" {
        return errors.New("export requires <in.pseudo> and <out.json>")
    }
    if err := gen.ExportToN8N(inPath, outPath); err != nil {
        return err
    }
    fmt.Printf("Exported: %s -> %s\n", inPath, outPath)
    return nil
}

func defaultWorkflowDir() string {
    return filepath.FromSlash("pkg/engine/workflows")
}

func compileAll() error {
    dir := defaultWorkflowDir()
    matches, err := filepath.Glob(filepath.Join(dir, "*.pseudo"))
    if err != nil {
        return err
    }
    if len(matches) == 0 {
        fmt.Println("No .pseudo files found in", dir)
        return nil
    }
    for _, in := range matches {
        base := filepath.Base(in)
        name := base[:len(base)-len(filepath.Ext(base))]
        out := filepath.Join(dir, "compiled", name+".go")
        if err := compileOne(in, out); err != nil {
            return err
        }
    }
    return nil
}

func exportAll() error {
    dir := defaultWorkflowDir()
    matches, err := filepath.Glob(filepath.Join(dir, "*.pseudo"))
    if err != nil {
        return err
    }
    if len(matches) == 0 {
        fmt.Println("No .pseudo files found in", dir)
        return nil
    }
    for _, in := range matches {
        base := filepath.Base(in)
        name := base[:len(base)-len(filepath.Ext(base))]
        out := filepath.Join(dir, "exports", name+".json")
        if err := exportOne(in, out); err != nil {
            return err
        }
    }
    return nil
}

func main() {
    if len(os.Args) < 2 {
        usage()
        os.Exit(2)
    }
    cmd := os.Args[1]
    switch cmd {
    case "compile":
        if len(os.Args) < 4 {
            usage()
            os.Exit(2)
        }
        if err := compileOne(os.Args[2], os.Args[3]); err != nil {
            fmt.Println("Error:", err)
            os.Exit(1)
        }
    case "export":
        if len(os.Args) < 4 {
            usage()
            os.Exit(2)
        }
        if err := exportOne(os.Args[2], os.Args[3]); err != nil {
            fmt.Println("Error:", err)
            os.Exit(1)
        }
    case "compile-all":
        if err := compileAll(); err != nil {
            fmt.Println("Error:", err)
            os.Exit(1)
        }
    case "export-all":
        if err := exportAll(); err != nil {
            fmt.Println("Error:", err)
            os.Exit(1)
        }
    default:
        usage()
        os.Exit(2)
    }
}

