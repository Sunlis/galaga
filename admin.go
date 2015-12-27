package edm

import (
  "io"
  "bytes"
  "net/http"
  "text/template"
  "util"
  "strconv"
  "encoding/json"
  "appengine"
  "appengine/datastore"
  // "appengine/user"
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
  Id             int    `datastore:"id" json:"id"`
  Name           string `datastore:"name" json:"name"`
  X              float64  `datastore:"x" json:"x"`
  Y              float64  `datastore:"y" json:"y"`
  Z              float64  `datastore:"z" json:"z"`
  Faction        string `datastore:"faction" json:"faction"`
  Population     int    `datastore:"pop" json:"population"`
  Government     string `datastore:"gov" json:"government"`
  Allegiance     string `datastore:"alleg" json:"allegiance"`
  State          string `datastore:"state" json:"state"`
  Security       string `datastore:"sec" json:"security"`
  PrimaryEconomy string `datastore:"primec" json:"primary_economy"`
  Power          string `datastore:"power" json:"power"`
  PowerState     string `datastore:"powerst" json:"power_state"`
  NeedsPermit    int    `datastore:"needperm" json:"needs_permit"`
  UpdatedAt      int    `datastore:"upd" json:"updated_at"`
  SimbadRef      string `datastore:"simbad" json:"simbad_ref"`
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

type Widget struct {
  Name string `datastore:"name"`
  Blob string `datastore:"blob"`
}


func handleAdmin(w http.ResponseWriter, r *http.Request) {
  ctx := appengine.NewContext(r)

  blob := ""
  check := r.FormValue("check")
  if check != "" {
    blob += "Ran query for " + check + "\n"
    q := datastore.NewQuery("System").Filter("name >=", check)
    for t := q.Run(ctx); ; {
      var s System
      _, err := t.Next(&s)
      if err == datastore.Done {
        break
      }
      if err != nil {
        break
      }
      blob += "Query returned: " + s.String() + "\n"
    }
  }

  if r.FormValue("widget_name") != "" {
    w_name := r.FormValue("widget_name")
    w_blob := r.FormValue("widget_blob")
    widget := Widget{
      Name: w_name,
      Blob: w_blob,
    }
    datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Widget", nil), &widget)
    blob += "Saved widget to datastore"
  }

  m, err := r.MultipartReader()
  if err == nil {
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
        blob += "FileName: " + part.FileName() + ", FormName: " + part.FormName() + ", Size: " + strconv.FormatInt(written, 10) + " b\n"

        var s *[]System
        err = json.Unmarshal([]byte(buf.String()), &s)
        util.CheckError("unmarshal", r, err)


        for k, v := range *s {
          key, err := datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "System", nil), &v)
          util.CheckError("database put", r, err)
          blob += "put: " + key.String() + " \n"
          blob += strconv.Itoa(k) + " " + v.String() + "\n"

          temp := new(System)
          err = datastore.Get(ctx, key, &temp)
          blob += "Wrote: " + v.String() + "\n"
        }
      }

      err = part.Close()
      util.CheckError("close part", r, err)
    }
  }


  t, err := template.ParseFiles("admin.html")
  util.CheckError("parse template", r, err)

  data := Admin {check, blob}

  err = t.Execute(w, data)
  util.CheckError("execute template", r, err)
}
