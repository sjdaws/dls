package auth

import (
    "encoding/json"
    "net/http"
    "time"

    "github.com/golang-jwt/jwt/v4"
    "github.com/sjdaws/dls/internal/global"
    "github.com/sjdaws/dls/internal/web"
)

type AuthCode struct {
    jwt.RegisteredClaims
    Challenge       string `json:"challenge"`
    KeyId           string `json:"kid"`
    KeyReference    string `json:"key_ref"`
    OriginReference string `json:"origin_ref"`
}

type Code struct {
    CodeChallenge   string `json:"code_challenge"`
    OriginReference string `json:"origin_ref"`
}

type CodeResponse struct {
    AuthCode      string   `json:"auth_code"`
    Prompts       []string `json:"prompts"`
    SyncTimestamp string   `json:"sync_timestamp"`
}

// CodeChallenge accepts a code and returns an authcode
func (a *Auth) CodeChallenge(response http.ResponseWriter, request *http.Request) {
    var body Code
    err := web.DecodeBody(request, &body)
    if err != nil {
        web.Error(request, response, "error unmarshaling json body", err, nil)
        return
    }

    currentTime := global.CurrentTime()
    expiryTime := currentTime.Add(15 * time.Minute)
    authCode, err := jwt.NewWithClaims(jwt.SigningMethodRS256, &AuthCode{
        jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expiryTime),
            IssuedAt:  jwt.NewNumericDate(currentTime),
        },
        body.CodeChallenge,
        "00000000-0000-0000-0000-000000000000",
        "00000000-0000-0000-0000-000000000000",
        body.OriginReference,
    }).SignedString(global.SigningKey)
    if err != nil {
        web.Error(request, response, "unable to marshal signed jwt", err, nil)
        return
    }

    reply, err := json.Marshal(&CodeResponse{
        AuthCode:      authCode,
        SyncTimestamp: currentTime.Format("2006-01-02T15:04:05.000000Z"),
    })
    if err != nil {
        web.Error(request, response, "unable to crease json response", err, nil)
        return
    }

    web.Respond(request, response, "application/json", "%s", reply)
}
