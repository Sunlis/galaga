package edm

import (
	// "io"
	"bytes"
  "strings"
	"net/http"
	"strconv"
	"text/template"
	"util"
  "mime/multipart"
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
	Blob  string
  SystemNames string
}

func (admin *Admin) Print(msg string) {
  admin.Blob += msg + "\n"
}

type Systems []struct {
	System System
}

type System struct {
	Id             int     `datastore:"id" json:"id"`
	Name           string  `datastore:"name" json:"name"`
	X              float64 `datastore:"x" json:"x"`
	Y              float64 `datastore:"y" json:"y"`
	Z              float64 `datastore:"z" json:"z"`
	Faction        string  `datastore:"faction" json:"faction"`
	Population     int     `datastore:"pop" json:"population"`
	Government     string  `datastore:"gov" json:"government"`
	Allegiance     string  `datastore:"alleg" json:"allegiance"`
	State          string  `datastore:"state" json:"state"`
	Security       string  `datastore:"sec" json:"security"`
	PrimaryEconomy string  `datastore:"primec" json:"primary_economy"`
	Power          string  `datastore:"power" json:"power"`
	PowerState     string  `datastore:"powerst" json:"power_state"`
	NeedsPermit    int     `datastore:"needperm" json:"needs_permit"`
	UpdatedAt      int     `datastore:"upd" json:"updated_at"`
	SimbadRef      string  `datastore:"simbad" json:"simbad_ref"`
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
	a := Admin{"", "", ""}

	err := r.ParseMultipartForm(128 * 1024 * 1024)
	if err == nil {
		if file, header, err := r.FormFile("systems"); err == nil {
			a.Print("Found file " + header.Filename)
			handleSystems(file, header, r, &a)
		}
		if value := r.FormValue("check"); value != "" {
      a.Check = value
      handleCheck(value, r, &a)
		} else {
			a.Print("No check")
		}
	} else {
		a.Print("No form data")
	}

  fetchSystemList(r, &a)

	t, err := template.ParseFiles("admin.html")
	util.CheckError("parse template", r, err)

	err = t.Execute(w, a)
	util.CheckError("execute template", r, err)
}

func fetchSystemList(r *http.Request, a *Admin) {
  ctx := appengine.NewContext(r)
  names := []string{}
  query := datastore.NewQuery("System").Project("id", "name").Order("name")
  for iter := query.Run(ctx); ; {
    var s System
    _, err := iter.Next(&s)
    if err != nil {
      break
    }
    names = append(names, s.Name)
  }
  if len(names) > 0 {
    a.SystemNames = "[\"" + strings.Join(names[:], "\", \"") + "\"]"
  } else {
    a.SystemNames = "[]"
  }
}

func handleCheck(check string, r *http.Request, a *Admin) {
  ctx := appengine.NewContext(r)
  a.Print("Searching datastore for \"" + check + "\"...")
  query := datastore.NewQuery("System").Filter("name =", check)
  results := 0
  for iter := query.Run(ctx); ; {
    var s System
    _, err := iter.Next(&s)
    if err == datastore.Done || err != nil {
      break
    }
    a.Print("Query found: " + s.String())
    results++
  }
  a.Print("Total results: " + strconv.Itoa(results))
}

func handleSystems(file multipart.File, header *multipart.FileHeader, r *http.Request, a *Admin) {
  ctx := appengine.NewContext(r)

  buf := new(bytes.Buffer)
  read, err := buf.ReadFrom(file)
  util.CheckError("read file", r, err)
  a.Print("Read " + strconv.FormatInt(read, 10) + " bytes from file")
  a.Print(buf.String())

	var s *[]System
	err = json.Unmarshal([]byte(buf.String()), &s)
	util.CheckError("unmarshal", r, err)
	for k, v := range *s {
		key, err := datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "System", nil), &v)
		util.CheckError("database put", r, err)
    a.Print("put: " + key.String())
		a.Print(strconv.Itoa(k) + " " + v.String())

		// temp := new(System)
		// err = datastore.Get(ctx, key, &temp)
		// a.Print("Wrote: " + v.String())
	}
}
