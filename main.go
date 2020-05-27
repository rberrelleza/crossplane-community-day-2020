package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/codegangsta/negroni"
	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

var db *pg.DB
var namespace string

// Guest holds the comment and an id
type Guest struct {
	ID      string `json:"-"`
	Comment string `json:"comment"`
}

func infoHandler(rw http.ResponseWriter, req *http.Request) {
	info := make(map[string]string)

	info["namespace"] = namespace
	info["pod"] = os.Getenv("HOSTNAME")
	info["golang"] = os.Getenv("GOLANG_VERSION")

	r, err := (json.MarshalIndent(info, "", "  "))
	if err != nil {
		fmt.Printf("error: %s\n", err)
		rw.WriteHeader(500)
		return
	}

	rw.Write(r)
}

func envHandler(rw http.ResponseWriter, req *http.Request) {
	environment := make(map[string]string)
	for _, item := range os.Environ() {
		splits := strings.Split(item, "=")
		key := splits[0]
		val := strings.Join(splits[1:], "=")
		if strings.Contains(item, "PGPASSWORD") {
			val = "**********"
		}

		environment[key] = val
	}

	r, err := (json.MarshalIndent(environment, "", "  "))
	if err != nil {
		fmt.Printf("error: %s\n", err)
		rw.WriteHeader(500)
		return
	}

	rw.Write(r)
}

func getHandler(rw http.ResponseWriter, req *http.Request) {
	var guests []Guest
	/*if err := db.Model(&guests).Select(); err != nil {
		fmt.Printf("error: %s\n", err)
		rw.WriteHeader(500)
		return
	}*/

	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(guests); err != nil {
		fmt.Printf("error: %s\n", err)
		rw.WriteHeader(500)
		return
	}
}

func submitHandler(rw http.ResponseWriter, req *http.Request) {
	g := Guest{
		ID:      uuid.New().String(),
		Comment: req.FormValue("comment"),
	}

	if g.Comment == "" {
		rw.WriteHeader(400)
		return
	}

	/*if err := db.Insert(&g); err != nil {
		fmt.Printf("error: %s\n", err)
		rw.WriteHeader(500)
		return
	}*/

	getHandler(rw, req)
}

func createSchema(db *pg.DB) error {
	return db.CreateTable(&Guest{}, &orm.CreateTableOptions{
		IfNotExists: true,
	})
}

func main() {

	db = pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("PGHOST"), os.Getenv("PGPORT")),
		User:     os.Getenv("PGUSER"),
		Password: os.Getenv("PGPASSWORD"),
		Database: os.Getenv("PGDATABASE"),
	})

	defer db.Close()

	err := createSchema(db)
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		namespace = string(b)
	}

	namespace = string(b)

	r := mux.NewRouter()
	r.Path("/env").Methods("GET").HandlerFunc(envHandler)
	r.Path("/info").Methods("GET").HandlerFunc(infoHandler)
	r.Path("/get").Methods("GET").HandlerFunc(getHandler)
	r.Path("/submit").Methods("POST").HandlerFunc(submitHandler)

	n := negroni.Classic()
	n.UseHandler(r)
	port := ":8080"
	if p, ok := os.LookupEnv("PORT"); ok {
		port = fmt.Sprintf(":%s", p)
	}
	n.Run(port)
}
