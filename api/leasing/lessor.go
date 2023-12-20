package leasing

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strings"
    "time"

    "github.com/google/uuid"
    "github.com/sjdaws/dls/api/auth"
    db "github.com/sjdaws/dls/database"
    "github.com/sjdaws/dls/internal/web"
)

type Application struct {
    FulfillmentContext     FulfillmentContext `json:"fulfillment_context"`
    LeaseProposalList      []LeaseProposal    `json:"lease_proposal_list"`
    ProposalEvaluationMode string             `json:"proposal_evaluation_mode"`
    ScopeReferenceList     []string           `json:"scope_ref_list"`
}

type ApplicationResponse struct {
    LeaseResultList []LeaseList `json:"lease_result_list"`
    Prompts         []string    `json:"prompts"`
    ResultCode      string      `json:"result_code"`
    SyncTimestamp   string      `json:"sync_timestamp"`
}

type FulfillmentContext struct {
    ClassReferenceList []string `json:"fulfillment_class_ref_list"`
}

type Lease struct {
    Created                 string `json:"created"`
    Expires                 string `json:"expires"`
    LicenceType             string `json:"licence_type"`
    OfflineLease            bool   `json:"offline_lease"`
    RecommendedLeaseRenewal string `json:"recommended_lease_renewal"`
    Reference               string `json:"ref"`
}

type LeaseList struct {
    Lease   Lease `json:"lease"`
    Ordinal int   `json:"ordinal"`
}

type ListLeasesResponse struct {
    ActiveLeaseList map[string]Lease `json:"active_lease_list"`
    Prompts         []string         `json:"prompts"`
    SyncTimestamp   string           `json:"sync_timestamp"`
}

type LeaseProposal struct {
    LicenceTypeQualifiers LicenceTypeQualifiers `json:"license_type_qualifiers"`
    Product               Product               `json:"product"`
}

type LicenceTypeQualifiers struct {
    Count int `json:"count"`
}

type Product struct {
    Name string `json:"product"`
}

type DeleteLeaseResponse struct {
    Prompts            []string         `json:"prompts"`
    ReleasedLeaseList  map[string]Lease `json:"released_lease_list"`
    ReleaseFailureList map[string]Lease `json:"release_failure_list"`
    SyncTimestamp      string           `json:"sync_timestamp"`
}

// CreateOriginLease creates a lease for a lessor
func (l *Leasing) CreateOriginLease(response http.ResponseWriter, request *http.Request) {
    var body Application
    err := web.DecodeBody(request, &body)
    if err != nil {
        web.Error(request, response, "error unmarshaling json body", err, nil)
        return
    }

    token, err := auth.ReadFromHeader(request)
    if err != nil {
        web.Error(request, response, "invalid jwt", err, &web.HttpError{Detail: "invalid token", StatusCode: http.StatusUnauthorized})
        return
    }

    currentTime := time.Now().UTC()
    expiryTime := currentTime.Add(90 * 24 * time.Hour)
    leaseList := make([]LeaseList, 0)
    for range body.ScopeReferenceList {
        reference := strings.ToUpper(fmt.Sprintf("%v", uuid.New()))

        err = l.database.CreateLease(&db.Lease{
            CreatedAt:       currentTime,
            ExpiresAt:       expiryTime,
            OriginReference: token.OriginReference,
            Reference:       reference,
        })
        if err != nil {
            web.Error(request, response, fmt.Sprintf("error creating lease in the database for origin '%s'", token.OriginReference), err, nil)
            return
        }

        leaseList = append(leaseList, LeaseList{
            Ordinal: 0,
            Lease: Lease{
                Created:                 currentTime.Format("2006-01-02T15:04:05.000000Z"),
                Expires:                 expiryTime.Format("2006-01-02T15:04:05.000000Z"),
                LicenceType:             "CONCURRENT_COUNTED_SINGLE",
                OfflineLease:            true,
                RecommendedLeaseRenewal: "0.15",
                Reference:               reference,
            },
        })
    }

    reply, err := json.Marshal(&ApplicationResponse{
        LeaseResultList: leaseList,
        ResultCode:      "SUCCESS",
        SyncTimestamp:   currentTime.Format("2006-01-02T15:04:05.000000Z"),
    })
    if err != nil {
        web.Error(request, response, "unable to marshal json response", err, nil)
        return
    }

    web.Respond(request, response, "application/json", "%s", reply)
}

// DeleteOriginLeases deletes leases for an origin reference
func (l *Leasing) DeleteOriginLeases(response http.ResponseWriter, request *http.Request) {
    token, err := auth.ReadFromHeader(request)
    if err != nil {
        web.Error(request, response, "invalid jwt", err, &web.HttpError{Detail: "invalid token", StatusCode: http.StatusUnauthorized})
        return
    }

    leases, err := l.database.GetLeases(token.OriginReference)
    if err != nil {
        web.Error(request, response, fmt.Sprintf("error fetching leases from the database for origin '%s'", token.OriginReference), err, nil)
        return
    }

    err = l.database.DeleteLeases(token.OriginReference)
    if err != nil {
        web.Error(request, response, fmt.Sprintf("error removing leases from the database for origin '%s'", token.OriginReference), err, nil)
        return
    }

    leaseList := map[string]Lease{}
    for _, lease := range leases {
        leaseList[lease.Reference] = Lease{
            Created:                 lease.CreatedAt.Format("2006-01-02T15:04:05.000000Z"),
            Expires:                 lease.ExpiresAt.Format("2006-01-02T15:04:05.000000Z"),
            LicenceType:             "CONCURRENT_COUNTED_SINGLE",
            OfflineLease:            true,
            RecommendedLeaseRenewal: "0.15",
            Reference:               lease.Reference,
        }
    }

    currentTime := time.Now().UTC()
    reply, err := json.Marshal(&DeleteLeaseResponse{
        ReleasedLeaseList: leaseList,
        SyncTimestamp:     currentTime.Format("2006-01-02T15:04:05.000000Z"),
    })
    if err != nil {
        web.Error(request, response, "unable to marshal json response", err, nil)
        return
    }

    web.Respond(request, response, "application/json", "%s", reply)
}

// GetOriginLeases fetches leases for an origin reference
func (l *Leasing) GetOriginLeases(response http.ResponseWriter, request *http.Request) {
    token, err := auth.ReadFromHeader(request)
    if err != nil {
        web.Error(request, response, "invalid jwt", err, &web.HttpError{Detail: "invalid token", StatusCode: http.StatusUnauthorized})
        return
    }

    leases, err := l.database.GetLeases(token.OriginReference)
    if err != nil {
        web.Error(request, response, fmt.Sprintf("error fetching leases from the database for origin '%s'", token.OriginReference), err, nil)
        return
    }

    leaseList := map[string]Lease{}
    for _, lease := range leases {
        leaseList[lease.Reference] = Lease{
            Created:                 lease.CreatedAt.Format("2006-01-02T15:04:05.000000Z"),
            Expires:                 lease.ExpiresAt.Format("2006-01-02T15:04:05.000000Z"),
            LicenceType:             "CONCURRENT_COUNTED_SINGLE",
            OfflineLease:            true,
            RecommendedLeaseRenewal: "0.15",
            Reference:               lease.Reference,
        }
    }

    currentTime := time.Now().UTC()
    reply, err := json.Marshal(&ListLeasesResponse{
        ActiveLeaseList: leaseList,
        SyncTimestamp:   currentTime.Format("2006-01-02T15:04:05.000000Z"),
    })
    if err != nil {
        web.Error(request, response, "unable to marshal json response", err, nil)
        return
    }

    web.Respond(request, response, "application/json", "%s", reply)
}
