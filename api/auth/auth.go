package auth

import (
    "net/http"
    "strings"

    "github.com/golang-jwt/jwt/v4"
    db "github.com/sjdaws/dls/database"
    "github.com/sjdaws/dls/internal/global"
)

type Auth struct {
    database *db.Database
}

type JWx struct {
    jwt.RegisteredClaims
    Challenge       string `json:"challenge,omitempty"`
    OriginReference string `json:"origin_ref"`
}

// New creates a new auth instance
func New(database *db.Database) *Auth {
    return &Auth{
        database: database,
    }
}

// ReadFromHeader returns the auth code from the authorization header
func ReadFromHeader(request *http.Request) (*JWx, error) {
    authCode := request.Header.Get("Authorization")
    if strings.EqualFold(authCode[0:6], "bearer") {
        authCode = authCode[7:]
    }

    return decode(authCode)
}

// decode decodes a JWx
func decode(jwx string) (*JWx, error) {
    token, err := jwt.ParseWithClaims(jwx, &JWx{}, func(token *jwt.Token) (interface{}, error) {
        return &global.SigningKey.PublicKey, nil
    })
    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*JWx); ok {
        return claims, nil
    }

    return nil, err
}
