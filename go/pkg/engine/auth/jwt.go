package auth

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
    "encoding/json"
    "errors"
    "strings"
    "time"
)

type Claims struct {
    Issuer   string `json:"iss,omitempty"`
    Subject  string `json:"sub,omitempty"`
    IssuedAt int64  `json:"iat,omitempty"`
    Expires  int64  `json:"exp,omitempty"`
}

func SignJWT(c Claims, secret string) (string, error) {
    if secret == "" { return "", errors.New("jwt: secret required") }
    header := map[string]string{"alg": "HS256", "typ": "JWT"}
    hb, _ := json.Marshal(header)
    cb, _ := json.Marshal(c)
    hEnc := b64url(hb)
    cEnc := b64url(cb)
    signingInput := hEnc + "." + cEnc
    sig := hmacSHA256(signingInput, secret)
    return signingInput + "." + b64url(sig), nil
}

func VerifyJWT(tok, secret string) (Claims, error) {
    var out Claims
    parts := strings.Split(tok, ".")
    if len(parts) != 3 { return out, errors.New("jwt: invalid token format") }
    signingInput := parts[0] + "." + parts[1]
    sig, err := b64urldec(parts[2])
    if err != nil { return out, errors.New("jwt: invalid signature encoding") }
    want := hmacSHA256(signingInput, secret)
    if !hmac.Equal(sig, want) { return out, errors.New("jwt: signature mismatch") }
    payload, err := b64urldec(parts[1])
    if err != nil { return out, errors.New("jwt: invalid payload encoding") }
    if err := json.Unmarshal(payload, &out); err != nil { return out, errors.New("jwt: invalid payload") }
    if out.Expires > 0 && time.Now().Unix() > out.Expires { return out, errors.New("jwt: expired") }
    return out, nil
}

func hmacSHA256(s, secret string) []byte {
    mac := hmac.New(sha256.New, []byte(secret))
    mac.Write([]byte(s))
    return mac.Sum(nil)
}

func b64url(b []byte) string {
    return base64.RawURLEncoding.EncodeToString(b)
}

func b64urldec(s string) ([]byte, error) {
    return base64.RawURLEncoding.DecodeString(s)
}

