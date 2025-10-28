package main

import (
    "bufio"
    "flag"
    "fmt"
    "io/fs"
    "log"
    "os"
    "path/filepath"
    "sort"
    "strings"

    _ "github.com/lib/pq"
    "github.com/bitesinbyte/ferret/pkg/db"
)

func main() {
    dir := flag.String("dir", "_data/_db", "directory with .sql migrations (applied in lexical order)")
    dsn := flag.String("database", os.Getenv("DATABASE_URL"), "Postgres DSN (or set DATABASE_URL)")
    dry := flag.Bool("dry", false, "print migration order without executing")
    flag.Parse()

    if *dsn == "" {
        log.Fatal("missing DATABASE_URL or --database")
    }
    files, err := listSQL(*dir)
    if err != nil { log.Fatal(err) }
    if len(files) == 0 { fmt.Println("no migrations found"); return }

    fmt.Println("migrations:")
    for _, f := range files { fmt.Println(" -", f) }
    if *dry { return }

    conn, err := db.Open(*dsn)
    if err != nil { log.Fatal(err) }
    defer conn.Close()

    for _, path := range files {
        b, err := os.ReadFile(path)
        if err != nil { log.Fatal(err) }
        sql := string(b)
        if strings.TrimSpace(sql) == "" { continue }
        fmt.Printf("\n-- applying %s\n", filepath.Base(path))
        if _, err := conn.Exec(sql); err != nil {
            // dump first 200 chars to help debug
            r := bufio.NewReader(strings.NewReader(sql))
            preview, _ := r.Peek(min(200, len(sql)))
            log.Fatalf("failed migration %s: %v\nSQL preview: %s", path, err, string(preview))
        }
    }
    fmt.Println("done")
}

func listSQL(dir string) ([]string, error) {
    var out []string
    err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
        if err != nil { return err }
        if d.IsDir() { return nil }
        if strings.HasSuffix(strings.ToLower(d.Name()), ".sql") {
            out = append(out, path)
        }
        return nil
    })
    if err != nil { return nil, err }
    sort.Strings(out)
    return out, nil
}

func min(a, b int) int { if a < b { return a }; return b }

