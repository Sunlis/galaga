package util

import (
  "os"
  "net/http"
  "appengine"
)

func CheckError(msg string, r *http.Request, err error) {
  if err != nil {
    ctx := appengine.NewContext(r)
    ctx.Errorf(msg + " - " + err.Error())
    os.Exit(1)
  }
}
