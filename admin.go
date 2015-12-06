package edm

import (
    "fmt"
    "net/http"
)

func init() {
    http.HandleFunc("/admin", handleAdmin)
}

func handleAdmin(w http.ResponseWriter, r *http.Request) {
    if (r.FormValue("foo") != "") {
      fmt.Fprint(w, r.FormValue("foo"))
    }
    http.ServeFile(w, r, "admin.html")
}
