package main

import (
    "fmt"
    "log"
    "github.com/bitesinbyte/ferret/pkg/auth"
)

// Small CLI demo to generate a session token and hash/verify it.
func main() {
    tok, err := auth.GenerateToken(32)
    if err != nil { log.Fatal(err) }
    hash := auth.HashToken(tok)
    fmt.Println("token:", tok)
    fmt.Println("hash:", hash)
    fmt.Println("verify (ok):", auth.CompareTokenHash(hash, tok))
    fmt.Println("verify (bad):", auth.CompareTokenHash(hash, tok+"x"))
}

