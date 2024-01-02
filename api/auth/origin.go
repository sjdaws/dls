package auth

import (
    "encoding/json"
    "net/http"

    "github.com/sjdaws/dls/internal/global"
    "github.com/sjdaws/dls/internal/web"
)

type Environment struct {
    ClientPlatformId        string         `json:"client_platform_id"`
    CpuSockets              int            `json:"cpu_sockets"`
    Fingerprint             Fingerprint    `json:"fingerprint"`
    GpuIdList               []string       `json:"gpu_id_list"`
    GuestDriverVersion      string         `json:"guest_driver_version"`
    HostDriverVersion       string         `json:"host_driver_version"`
    Hostname                string         `json:"hostname"`
    HvPlatform              string         `json:"hv_platform"`
    IpAddressList           []string       `json:"ip_address_list"`
    OperatingSystemPlatform string         `json:"os_platform"`
    OperatingSystemVersion  string         `json:"os_version"`
    PhysicalCores           int            `json:"physical_cores"`
    RawEnvironment          RawEnvironment `json:"raw_env"`
}

type Fingerprint struct {
    MacAddressList []string `json:"mac_address_list"`
}

type Origin struct {
    CandidateOriginReference string      `json:"candidate_origin_ref"`
    Environment              Environment `json:"environment"`
}

type OriginResponse struct {
    Environment        Environment `json:"environment"`
    NodeQueryOrder     []string    `json:"node_query_order"`
    NodeUrlList        []string    `json:"node_url_list"`
    OriginReference    string      `json:"origin_ref"`
    Prompts            []string    `json:"prompts"`
    ServicePortSetList []string    `json:"svc_port_set_list"`
    SyncTimestamp      string      `json:"sync_timestamp"`
}

type RawEnvironment struct {
    ClientPlatformId        string      `json:"client_platform_id"`
    CpuSockets              int         `json:"cpu_sockets"`
    Fingerprint             Fingerprint `json:"fingerprint"`
    GpuIdList               []string    `json:"gpu_id_list"`
    GuestDriverVersion      string      `json:"guest_driver_version"`
    HostDriverVersion       string      `json:"host_driver_version"`
    Hostname                string      `json:"hostname"`
    HvPlatform              string      `json:"hv_platform"`
    IpAddressList           []string    `json:"ip_address_list"`
    OperatingSystemPlatform string      `json:"os_platform"`
    OperatingSystemVersion  string      `json:"os_version"`
    PhysicalCores           int         `json:"physical_cores"`
}

// FindOrCreateOrigin finds an existing origin, or creates a new one in the database
func (a *Auth) FindOrCreateOrigin(response http.ResponseWriter, request *http.Request) {
    var body Origin
    err := web.DecodeBody(request, &body)
    if err != nil {
        web.Error(request, response, "error unmarshaling json body", err, nil)
        return
    }

    body.Environment.RawEnvironment = RawEnvironment{
        ClientPlatformId:        body.Environment.ClientPlatformId,
        CpuSockets:              body.Environment.CpuSockets,
        Fingerprint:             body.Environment.Fingerprint,
        GuestDriverVersion:      body.Environment.GuestDriverVersion,
        GpuIdList:               body.Environment.GpuIdList,
        HostDriverVersion:       body.Environment.HostDriverVersion,
        Hostname:                body.Environment.Hostname,
        HvPlatform:              body.Environment.HvPlatform,
        IpAddressList:           body.Environment.IpAddressList,
        OperatingSystemPlatform: body.Environment.OperatingSystemPlatform,
        OperatingSystemVersion:  body.Environment.OperatingSystemVersion,
        PhysicalCores:           body.Environment.PhysicalCores,
    }
    body.Environment.Fingerprint.MacAddressList = []string{}

    reply, err := json.Marshal(&OriginResponse{
        Environment:     body.Environment,
        OriginReference: body.CandidateOriginReference,
        SyncTimestamp:   global.CurrentTime().Format("2006-01-02T15:04:05.000000Z"),
    })
    if err != nil {
        web.Error(request, response, "unable to marshal json response", err, nil)
        return
    }

    web.Respond(request, response, "application/json", "%s", reply)
}
