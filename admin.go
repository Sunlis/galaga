package edm

import (
	// "io"
  "fmt"
	"bytes"
  "strings"
	"net/http"
	"strconv"
	"text/template"
  "util"
	"time"
  "mime/multipart"
	"encoding/json"
	"appengine"
  "appengine/datastore"
	"appengine/search"
	// "appengine/user"
)

func init() {
	http.HandleFunc("/admin", handleAdmin)
}

type Data struct {
  Updated int `datastore:"updated"`
}

type Admin struct {
	Check string
	Blob  string
}

func (admin *Admin) Print(msg string) {
  admin.Blob += msg + "\n"
}

func (admin *Admin) Printf(msg string, parts ...interface{}) {
  admin.Print(fmt.Sprintf(msg, parts...))
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

type SearchableSystem struct {
  Id       float64
  Name     string
  RealName string
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

func getData(r *http.Request) Data {
  ctx := appengine.NewContext(r)
  query := datastore.NewQuery("Data")
  for iter := query.Run(ctx); ; {
    var d Data
    _, err := iter.Next(&d)
    if err != nil {
      break
    }
    return d
  }
  return Data{0}
}

func handleAdmin(w http.ResponseWriter, r *http.Request) {
	a := Admin{"", ""}

  data := getData(r)

	err := r.ParseMultipartForm(128 * 1024 * 1024)
	if err == nil {
		if file, header, err := r.FormFile("systems"); err == nil {
			handleSystems(file, header, r, &a, &data)
		}
		if value := r.FormValue("check"); value != "" {
      a.Check = value
      handleCheck(value, r, &a)
		}
	}

	t, err := template.ParseFiles("templates/admin.html")
	util.CheckError("parse template", r, err)

	err = t.Execute(w, a)
	util.CheckError("execute template", r, err)
}

func handleCheck(check string, r *http.Request, a *Admin) {
  ctx := appengine.NewContext(r)
  a.Print("Searching for \"" + check + "\"...")
  index, err := search.Open("system_names")
  util.CheckError("open index", r, err)

  for t := index.Search(ctx, util.SpaceOut(check), nil); ; {
    var sys SearchableSystem
    _, err = t.Next(&sys)
    if err != nil {
      break
    }

    if checkSystem(check, sys.RealName) {
      a.Printf("Found: %s (ID: %d)", sys.RealName, int(sys.Id))
    }

    // query := datastore.NewQuery("System").Filter("id =", int(sys.Id))
    // for iter := query.Run(ctx); ; {
    //   var s System
    //   _, err := iter.Next(&s)
    //   if err == datastore.Done || err != nil {
    //     break
    //   }
    //   if checkSystem(check, s.Name) {
    //     a.Printf("System { Name: \"%s\", Id: %d }", s.Name, s.Id)
    //   }
    // }
  }
}

func handleSystems(file multipart.File, header *multipart.FileHeader, r *http.Request, a *Admin, data *Data) {
  ctx := appengine.NewContext(r)

  buf := new(bytes.Buffer)
  _, err := buf.ReadFrom(file)
  util.CheckError("read file", r, err)

	var s []System
	err = json.Unmarshal([]byte(buf.String()), &s)
	util.CheckError("unmarshal", r, err)

  index, err := search.Open("system_names")
  util.CheckError("open index", r, err)

  addCount := 0
  skipCount := 0

  var keys []*datastore.Key

  var toSave []System

  for _, v := range s {
    simple := &SearchableSystem{
      Id: float64(v.Id),
      Name: util.SpaceOut(v.Name),
      RealName: v.Name,
    }
    id := fmt.Sprintf("sys-%08d", simple.Id)
    err = index.Delete(ctx, id)
    util.CheckError("remove document", r, err)
    _, err = index.Put(ctx, id, simple)
    util.CheckError("put document", r, err)

    if r.FormValue("datecheck") == "" || v.UpdatedAt > data.Updated {
      keys = append(keys, datastore.NewIncompleteKey(ctx, "System", nil))
      query := datastore.NewQuery("System").Filter("id =", v.Id)
      for iter := query.Run(ctx); ; {
        var s System
        key, err := iter.Next(&s)
        if err == datastore.Done || err != nil {
          break
        }
        err = datastore.Delete(ctx, key)
        util.CheckError("remove old system", r, err)
        addCount++
      }
      toSave = append(toSave, v)
    } else {
      skipCount++
    }
  }

  _, err = datastore.PutMulti(ctx, keys, toSave)
  util.CheckError("putmulti", r, err)

  a.Printf("Removed %d. Skipped %d. Added %d.", addCount, skipCount, len(toSave))

  updateData(r, data)
}

func updateData(r *http.Request, data *Data) {
  ctx := appengine.NewContext(r)
  query := datastore.NewQuery("Data")
  for iter := query.Run(ctx); ; {
    var d Data
    key, err := iter.Next(&d)
    if err != nil {
      break
    }
    datastore.Delete(ctx, key)
    util.CheckError("remove old data entity", r, err)
  }

  data.Updated = int(time.Now().Unix())
  datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Data", nil), data)
}

func checkSystem(query string, name string) bool {
  return strings.Contains(strings.ToLower(name), strings.ToLower(query))
}
