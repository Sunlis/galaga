package util

import (
  "os"
  "net/http"
  "appengine"
  "strings"
)

func CheckError(msg string, r *http.Request, err error) {
  if err != nil {
    ctx := appengine.NewContext(r)
    ctx.Errorf(msg + " - " + err.Error())
    os.Exit(1)
  }
}

func SpaceOut(str string) string {
  return strings.ToLower(strings.Join(strings.Split(strings.Replace(str, " ", "", -1), ""), " "))
}
