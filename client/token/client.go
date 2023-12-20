package token

import (
    "crypto/x509"
    "encoding/pem"
    "fmt"
    "net/http"
    "strings"

    "github.com/golang-jwt/jwt/v4"
    "github.com/google/uuid"
    "github.com/sjdaws/dls/internal/global"
    "github.com/sjdaws/dls/internal/web"
)

type ClientToken struct {
    jwt.RegisteredClaims
    FulfillmentClassReferenceList         []string                              `json:"fulfillment_class_ref_list"`
    ScopeReferenceList                    []string                              `json:"scope_ref_list"`
    ServiceInstanceConfiguration          ServiceInstanceConfiguration          `json:"service_instance_configuration"`
    ServiceInstancePublicKeyConfiguration ServiceInstancePublicKeyConfiguration `json:"service_instance_public_key_configuration"`
    UpdateMode                            string                                `json:"update_mode"`
}

type NodeUrlList struct {
    Index               int    `json:"idx"`
    RestUrl             string `json:"url_qr"`
    ServicePortSetIndex int    `json:"svc_port_set_idx"`
    Url                 string `json:"url"`
}

type ServiceInstanceConfiguration struct {
    NodeUrlList                           []NodeUrlList        `json:"node_url_list"`
    NvidiaLicenceServiceInstanceReference string               `json:"nls_service_instance_ref"`
    ServicePortSetList                    []ServicePortSetList `json:"svc_port_set_list"`
}

type ServiceInstancePublicKeyConfiguration struct {
    KeyRetentionMode                       string                                 `json:"key_retention_mode"`
    ServiceInstancePublicKeyModuloExponent ServiceInstancePublicKeyModuloExponent `json:"service_instance_public_key_me"`
    ServiceInstancePublicKeyPem            string                                 `json:"service_instance_public_key_pem"`
}

type ServiceInstancePublicKeyModuloExponent struct {
    Exponent int    `json:"exp"`
    Modulo   string `json:"mod"`
}

// ServicePortMap must have service name first as nvidia parses this using an index rather than key
type ServicePortMap struct {
    Service string `json:"service"`
    Port    int    `json:"port"`
}

type ServicePortSetList struct {
    DelegationName string           `json:"d_name"`
    Index          int              `json:"idx"`
    ServicePortMap []ServicePortMap `json:"svc_port_map"`
}

func Download(response http.ResponseWriter, request *http.Request) {
    privateKey := global.SigningKey
    publicKeyDer, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
    if err != nil {
        web.Error(request, response, "unable to marshal public key", err, nil)
        return
    }

    pubKeyBlock := pem.Block{
        Type:    "PUBLIC KEY",
        Headers: nil,
        Bytes:   publicKeyDer,
    }
    publicKey := string(pem.EncodeToMemory(&pubKeyBlock))

    currentTime := global.CurrentTime()
    expiryTime := currentTime.AddDate(12, 0, 0)

    token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, &ClientToken{
        jwt.RegisteredClaims{
            Audience:  []string{"NLS Licensed Client"},
            ExpiresAt: jwt.NewNumericDate(expiryTime),
            ID:        strings.ToUpper(fmt.Sprintf("%v", uuid.New())),
            IssuedAt:  jwt.NewNumericDate(currentTime),
            Issuer:    "NLS Service Instance",
            NotBefore: jwt.NewNumericDate(currentTime),
        },
        []string{},
        []string{global.ScopeReference},
        ServiceInstanceConfiguration{
            NodeUrlList: []NodeUrlList{
                {
                    Index:               0,
                    RestUrl:             global.HttpHost,
                    ServicePortSetIndex: 0,
                    Url:                 global.HttpHost,
                },
            },
            NvidiaLicenceServiceInstanceReference: global.InstanceReference,
            ServicePortSetList: []ServicePortSetList{
                {
                    DelegationName: "DLS",
                    Index:          0,
                    ServicePortMap: []ServicePortMap{
                        {
                            Service: "auth",
                            Port:    global.HttpPort,
                        },
                        {
                            Service: "lease",
                            Port:    global.HttpPort,
                        },
                    },
                },
            },
        },
        ServiceInstancePublicKeyConfiguration{
            KeyRetentionMode: "LATEST_ONLY",
            ServiceInstancePublicKeyModuloExponent: ServiceInstancePublicKeyModuloExponent{
                Exponent: privateKey.PublicKey.E,
                Modulo:   fmt.Sprintf("%x", privateKey.PublicKey.N),
            },
            ServiceInstancePublicKeyPem: strings.TrimSpace(publicKey),
        },
        "ABSOLUTE",
    }).SignedString(global.SigningKey)
    if err != nil {
        web.Error(request, response, "unable to marshal signed jwt", err, nil)
        return
    }

    response.Header().Add("Content-Disposition", "attachment; filename=client_configuration_token.tok")
    web.Respond(request, response, "text/plain", "%s", token)
}
