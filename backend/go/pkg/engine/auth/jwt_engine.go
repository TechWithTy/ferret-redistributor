package auth

import "fmt"

type JWTAuth struct{}

func (a JWTAuth) Authenticate() {
    fmt.Println("AuthEngine: Authenticating via JWT tokens")
}

