package edm

import (
    "fmt"
    "net/http"
)

func init() {
    http.HandleFunc("/system?name=*", findSystem)
}

func findSystem(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, r.URL.Query()["name"])
}
