package auth

import "fmt"

type OAuth struct{}

func (o OAuth) AuthenticateProvider(provider string) {
    fmt.Printf("AuthEngine: Authenticating via OAuth provider: %s\n", provider)
}

