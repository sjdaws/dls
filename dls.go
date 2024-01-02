package main

import (
    "errors"
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/gorilla/mux"
    authv1 "github.com/sjdaws/dls/api/auth"
    leasingv1 "github.com/sjdaws/dls/api/leasing"
    "github.com/sjdaws/dls/client/health"
    "github.com/sjdaws/dls/client/token"
    "github.com/sjdaws/dls/internal/global"
    "github.com/sjdaws/dls/internal/web"
)

func main() {
    auth := authv1.New()
    leasing := leasingv1.New()

    router := mux.NewRouter()
    router.HandleFunc("/auth/v1/code", auth.CodeChallenge).Methods("POST")
    router.HandleFunc("/auth/v1/origin", auth.FindOrCreateOrigin).Methods("POST")
    router.HandleFunc("/auth/v1/origin/update", auth.FindOrCreateOrigin).Methods("POST")
    router.HandleFunc("/auth/v1/token", auth.TokenExchange).Methods("POST")
    router.HandleFunc("/health", health.GetHealth).Methods("GET")
    router.HandleFunc("/leasing/v1/lease/{id}", leasing.DeleteLease).Methods("DELETE")
    router.HandleFunc("/leasing/v1/lease/{id}", leasing.UpdateLease).Methods("PUT")
    router.HandleFunc("/leasing/v1/lessor", leasing.CreateOriginLease).Methods("POST")
    router.HandleFunc("/leasing/v1/lessor/leases", leasing.DeleteOriginLeases).Methods("DELETE")
    router.HandleFunc("/leasing/v1/lessor/leases", leasing.GetOriginLeases).Methods("GET")
    router.HandleFunc("/token/download", token.Download).Methods("GET")
    router.MethodNotAllowedHandler = http.HandlerFunc(invalidRoute)
    router.NotFoundHandler = http.HandlerFunc(invalidRoute)
    http.Handle("/", router)

    errs := make(chan error, 1)
    go serveHTTP(router, errs)
    log.Fatal(<-errs)
}

func invalidRoute(response http.ResponseWriter, request *http.Request) {
    web.Error(request, response, "route requested but not found", errors.New("invalid route"), nil)
}

func serveHTTP(router *mux.Router, errs chan<- error) {
    server := &http.Server{
        Addr:         fmt.Sprintf(":%d", global.ContainerPort),
        Handler:      router,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
    }
    log.Printf("Listening for HTTP connections on port %d", global.ContainerPort)
    errs <- server.ListenAndServe()
}
