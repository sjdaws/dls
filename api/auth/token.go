package auth

import (
    "crypto/sha256"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "net/http"
    "strings"
    "time"

    "github.com/golang-jwt/jwt/v4"
    "github.com/sjdaws/dls/internal/global"
    "github.com/sjdaws/dls/internal/web"
)

type AuthToken struct {
    jwt.RegisteredClaims
    KeyId           string `json:"kid"`
    KeyReference    string `json:"key_ref"`
    OriginReference string `json:"origin_ref"`
}

type Token struct {
    AuthCode     string `json:"auth_code"`
    CodeVerifier string `json:"code_verifier"`
}

type TokenError struct {
    Detail     string `json:"detail"`
    StatusCode int    `json:"status"`
}

type TokenResponse struct {
    AuthToken     string `json:"auth_token"`
    Expires       string `json:"expires"`
    SyncTimestamp string `json:"sync_timestamp"`
}

// TokenExchange exchanges an auth token for a token response
func (a *Auth) TokenExchange(response http.ResponseWriter, request *http.Request) {
    var body Token
    err := web.DecodeBody(request, &body)
    if err != nil {
        web.Error(request, response, "error unmarshaling json body", err, nil)
        return
    }

    token, err := decode(body.AuthCode)
    if err != nil {
        web.Error(request, response, "unable to decode jwx", err, nil)
        return
    }

    // Verify challenge
    sha := sha256.New()
    sha.Write([]byte(body.CodeVerifier))
    if strings.TrimRight(base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s", sha.Sum(nil)))), "=") != token.Challenge {
        web.Error(request, response, "unable to verify challenge", err, &web.HttpError{
            Detail:     "challenge did not match code verifier",
            StatusCode: http.StatusUnauthorized,
        })
        return
    }

    currentTime := global.CurrentTime()
    expiryTime := currentTime.Add(24 * time.Hour)
    authToken, err := jwt.NewWithClaims(jwt.SigningMethodRS256, &AuthToken{
        jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expiryTime),
            IssuedAt:  jwt.NewNumericDate(currentTime),
            Issuer:    "https://cls.nvidia.org",
            NotBefore: jwt.NewNumericDate(currentTime),
        },
        "00000000-0000-0000-0000-000000000000",
        "00000000-0000-0000-0000-000000000000",
        token.OriginReference,
    }).SignedString(global.SigningKey)
    if err != nil {
        web.Error(request, response, "unable to marshal signed jwt", err, nil)
        return
    }

    reply, err := json.Marshal(&TokenResponse{
        AuthToken:     authToken,
        Expires:       expiryTime.Format("2006-01-02T15:04:05.000000Z"),
        SyncTimestamp: currentTime.Format("2006-01-02T15:04:05.000000Z"),
    })
    if err != nil {
        web.Error(request, response, "unable to crease json response", err, nil)
        return
    }

    web.Respond(request, response, "application/json", "%s", reply)
}
