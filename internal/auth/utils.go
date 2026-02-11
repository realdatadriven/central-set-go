package auth

import (
    // ...
    "crypto/sha256"
    "encoding/base64"
    "crypto/rand"
    "strings"

    "golang.org/x/oauth2"
)

// GeneratePKCEVerifier creates a high-entropy code_verifier (43–128 chars recommended)
func GeneratePKCEVerifier() (string, error) {
    // 32 bytes → 43 chars base64url after encoding
    b := make([]byte, 32)
    if _, err := rand.Read(b); err != nil {
        return "", err
    }

    // base64url encoding without padding
    v := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b)

    // RFC 7636: length between 43 and 128 chars
    if len(v) < 43 || len(v) > 128 {
        return "", errors.New("generated verifier length out of range")
    }

    return v, nil
}

// S256CodeChallenge creates code_challenge = BASE64URL(SHA256(code_verifier))
func S256CodeChallenge(verifier string) string {
    h := sha256.Sum256([]byte(verifier))
    return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(h[:])
}