package global

import (
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "log"
    "os"
)

var (
    SigningKey *rsa.PrivateKey
)

// loadSigningKey validates and loads signing key
func loadSigningKey() {
    var block *pem.Block

    if signingKey != "" {
        block, _ = pem.Decode([]byte(signingKey))
    }

    if signingKeyPath != "" {
        contents, err := os.ReadFile(signingKeyPath)
        if err != nil {
            log.Fatalf("unable to read signing key located at %s: %v", signingKeyPath, err)
        }

        block, _ = pem.Decode(contents)
    }

    privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
    if err != nil {
        log.Fatalf("unable to parse signing key located at %s: %v", signingKeyPath, err)
    }

    SigningKey = privateKey.(*rsa.PrivateKey)
}
