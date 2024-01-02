package global

import (
    "log"
    "os"
    "time"
)

var (
    ContainerPort       = aToI(os.Getenv("CONTAINER_PORT"))
    Debug               = aToB(os.Getenv("DEBUG"))
    HttpHost            = os.Getenv("FORWARDER_HOST")
    HttpPort            = aToI(os.Getenv("FORWARDER_PORT"))
    InstanceReference   = os.Getenv("INSTANCE_REFERENCE")
    KeyReference        = os.Getenv("KEY_REFERENCE")
    LeaseDuration       = aToD(os.Getenv("LEASE_DURATION"))
    LeaseRenewalPercent = aToP(os.Getenv("LEASE_RENEWAL_PERCENT"))
    NotificationUrls    = os.Getenv("NOTIFICATION_URLS")
    ScopeReference      = os.Getenv("SCOPE_REFERENCE")
    signingKeyPath      = os.Getenv("SIGNING_KEY_PATH")
)

func setEnvDefaults() {
    if HttpHost == "" {
        HttpHost = "localhost"
    }

    if HttpPort == 0 {
        HttpPort = 80
    }

    if ContainerPort == 0 {
        ContainerPort = HttpPort
    }

    if InstanceReference == "" {
        InstanceReference = "00000000-0000-0000-0000-000000000000"
    }

    if KeyReference == "" {
        KeyReference = "10000000-0000-0000-0000-000000000001"
    }

    if LeaseDuration == 0 {
        LeaseDuration = 90 * 24 * time.Hour
    }

    if LeaseRenewalPercent == 0 {
        LeaseRenewalPercent = 0.15
    }

    if ScopeReference == "" {
        ScopeReference = "20000000-0000-0000-0000-000000000002"
    }

    if signingKeyPath == "" {
        log.Fatal("The SIGNING_KEY_PATH environment variable is mandatory")
    }
}
