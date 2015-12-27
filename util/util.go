package util

import (
  "os"
  "net/http"
  "appengine"
)

func CheckError(msg string, r *http.Request, err error) {
  if err != nil {
    c := appengine.NewContext(r)
    c.Infof(msg + " - " + err.Error())
    os.Exit(1)
  }
}
