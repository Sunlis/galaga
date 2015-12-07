package edm

import (
  "io"
  "bytes"
  "net/http"
  "text/template"
  "util"
  "strconv"
)

func init() {
  http.HandleFunc("/admin", handleAdmin)
}

type Admin struct {
  Check string
  Blob string
}

func handleAdmin(w http.ResponseWriter, r *http.Request) {
  m, err := r.MultipartReader()
  util.CheckError("get multipart reader", r, err)
  for {
    part, err := m.NextPart()
    if err == io.EOF {
      break
    } else if err != nil {
      util.CheckError("next form part", r, err)
    }


    buf := new(bytes.Buffer)
    written, err := io.Copy(buf, part)
    io.WriteString(w, "FileName: " + part.FileName() + ", FormName: " + part.FormName() + ", Size: " + strconv.FormatInt(written, 10) + " b\n")
    io.WriteString(w, buf.String())

    // ordinarily, you would check if this is the file you were interested in, and stream it
    // to a backing store or to disk if it was. Note that you should probably wrap part with an
    // io.LimitReader to avoid DOS attacks
    //
    // For examples sake, we skip this file by Closing it right away.

    err = part.Close()
    util.CheckError("close part", r, err)
  }
  return

  t, err := template.ParseFiles("admin.html")
  util.CheckError("parse template", r, err)

  blob := ""

  // if r.FormValue("systems") != "" {
  //   file, header, err := r.FormFile("systems")
  //   util.CheckError("form file", r, err)
  //   _ = file
  //   _ = header

  //   buf := new(bytes.Buffer)
  //   buf.ReadFrom(file)

  //   blob = "Got file: " + r.FormValue("systems")
  //   // blob += " - " + buf.String()
  // }

  data := Admin {r.FormValue("check"), blob}

  err = t.Execute(w, data)
  util.CheckError("execute template", r, err)
}
