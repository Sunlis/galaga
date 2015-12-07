package edm

import (
  "io"
  "bytes"
  "net/http"
  "text/template"
  "util"
  "strconv"
  "encoding/json"
)

func init() {
  http.HandleFunc("/admin", handleAdmin)
}

type Admin struct {
  Check string
  Blob string
}

type Systems []struct {
  System System
}

type System struct {
  Id             int    `json:"id"`
  Name           string `json:"name"`
  X              float64  `json:"x"`
  Y              float64  `json:"y"`
  Z              float64  `json:"z"`
  Faction        string `json:"faction"`
  Population     int    `json:"population"`
  Government     string `json:"government"`
  Allegiance     string `json:"allegiance"`
  State          string `json:"state"`
  Security       string `json:"security"`
  PrimaryEconomy string `json:"primary_economy"`
  Power          string `json:"power"`
  PowerState     string `json:"power_state"`
  NeedsPermit    int    `json:"needs_permit"`
  UpdatedAt      int    `json:"updated_at"`
  SimbadRef      string `json:"simbad_ref"`
}

func (s System) String() string {
  return "{\n" +
    "  Id: " + strconv.Itoa(s.Id) + ",\n" +
    "  Name: \"" + s.Name + "\",\n" +
    "  X: " + strconv.FormatFloat(s.X, 'f', 6, 64) + ",\n" +
    "  Y: " + strconv.FormatFloat(s.Y, 'f', 6, 64) + ",\n" +
    "  Z: " + strconv.FormatFloat(s.Z, 'f', 6, 64) + ",\n" +
    "  Faction: \"" + s.Faction + "\",\n" +
    "  Population: " + strconv.Itoa(s.Population) + ",\n" +
    "  Government: \"" + s.Government + "\",\n" +
    "  Allegiance: \"" + s.Allegiance + "\",\n" +
    "  State: \"" + s.State + "\",\n" +
    "  Security: \"" + s.Security + "\",\n" +
    "  PrimaryEconomy: \"" + s.PrimaryEconomy + "\",\n" +
    "  Power: \"" + s.Power + "\",\n" +
    "  PowerState: \"" + s.PowerState + "\",\n" +
    "  NeedsPermit: " + strconv.Itoa(s.NeedsPermit) + ",\n" +
    "  UpdatedAt: " + strconv.Itoa(s.UpdatedAt) + ",\n" +
    "  SimbadRef: \"" + s.SimbadRef + "\"\n}"
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

    if part.FormName() == "systems" {
      buf := new(bytes.Buffer)
      written, err := io.Copy(buf, part)
      io.WriteString(w, "FileName: " + part.FileName() + ", FormName: " + part.FormName() + ", Size: " + strconv.FormatInt(written, 10) + " b\n")

      var s *[]System
      err = json.Unmarshal([]byte(buf.String()), &s)
      util.CheckError("unmarshal", r, err)

      io.WriteString(w, strconv.Itoa(len(*s)) + "\n")

      for k, v := range *s {
        io.WriteString(w, strconv.Itoa(k) + " " + v.String() + "\n")
      }
    }

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
