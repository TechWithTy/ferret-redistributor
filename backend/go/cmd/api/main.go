package main

import (
    "fmt"
    "log"
    "os"

    "github.com/bitesinbyte/ferret/pkg/api/server"
    _ "github.com/lib/pq"
)

func main() {
    port := os.Getenv("API_PORT")
    if port == "" { port = "8080" }
    s := server.New()
    addr := ":" + port
    log.Printf("api listening on %s", addr)
    if err := s.Run(addr); err != nil {
        fmt.Println("api server error:", err)
        os.Exit(1)
    }
}
