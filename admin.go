package edm

import (
  "net/http"
  "text/template"
  "os"
)

func init() {
  http.HandleFunc("/admin", handleAdmin)
}

// type Admin struct {
//   Check string
//   Blob string
// }

const test = `
This should work
`

func handleAdmin(w http.ResponseWriter, r *http.Request) {
  t, err := template.ParseFiles("admin.html")
  if err != nil {
    os.Exit(1)
  }

  // data := Admin {}//r.FormValue("check"), r.FormValue("blob")}

  err = t.Execute(w, nil)
  if err != nil {
    os.Exit(1)
  }
}
