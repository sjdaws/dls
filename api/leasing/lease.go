package leasing

import (
    "encoding/json"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/sjdaws/dls/api/auth"
    "github.com/sjdaws/dls/internal/global"
    "github.com/sjdaws/dls/internal/web"
)

type LeaseResponse struct {
    ExpiresAt               string   `json:"expires,omitempty"`
    OfflineLease            bool     `json:"offline_lease,omitempty"`
    Prompts                 []string `json:"prompts"`
    RecommendedLeaseRenewal float32  `json:"recommended_lease_renewal,omitempty"`
    Reference               string   `json:"lease_ref"`
    SyncTimestamp           string   `json:"sync_timestamp"`
}

// DeleteLease removes a lease from the database
func (l *Leasing) DeleteLease(response http.ResponseWriter, request *http.Request) {
    parameters := mux.Vars(request)
    leaseId := parameters["id"]

    _, err := auth.ReadFromHeader(request)
    if err != nil {
        web.Error(request, response, "invalid jwt", err, &web.HttpError{Detail: "invalid token", StatusCode: http.StatusUnauthorized})
        return
    }

    currentTime := global.CurrentTime()
    reply, err := json.Marshal(&LeaseResponse{
        Reference:     leaseId,
        SyncTimestamp: currentTime.Format("2006-01-02T15:04:05.000000Z"),
    })
    if err != nil {
        web.Error(request, response, "unable to marshal json response", err, nil)
        return
    }

    web.Respond(request, response, "application/json", "%s", reply)
}

// UpdateLease updates an existing lease
func (l *Leasing) UpdateLease(response http.ResponseWriter, request *http.Request) {
    parameters := mux.Vars(request)
    leaseId := parameters["id"]

    _, err := auth.ReadFromHeader(request)
    if err != nil {
        web.Error(request, response, "invalid jwt", err, &web.HttpError{Detail: "invalid token", StatusCode: http.StatusUnauthorized})
        return
    }

    currentTime := global.CurrentTime()
    expiryTime := currentTime.Add(global.LeaseDuration)

    reply, err := json.Marshal(&LeaseResponse{
        ExpiresAt:               expiryTime.Format("2006-01-02T15:04:05.000000Z"),
        OfflineLease:            true,
        RecommendedLeaseRenewal: global.LeaseRenewalPercent,
        Reference:               leaseId,
        SyncTimestamp:           currentTime.Format("2006-01-02T15:04:05.000000Z"),
    })
    if err != nil {
        web.Error(request, response, "unable to marshal json response", err, nil)
        return
    }

    web.Respond(request, response, "application/json", "%s", reply)
}
