package main

import (
	"log"
	"context"
	"net/http"

	"github.com/apuigsech/rest-layer-sql"

	"github.com/rs/rest-layer/resource"
	"github.com/rs/rest-layer/rest"
	"github.com/rs/rest-layer/schema"

	_ "github.com/mattn/go-sqlite3"
)

const (
	DB_DRIVER		= "sqlite3"
	DB_SOURCE		= "file::memory:?cache=shared"
	DB_TABLE 		= "units"

	DB_TABLE_UP		= "CREATE TABLE IF NOT EXISTS units (id VARCHAR(128) PRIMARY KEY,etag VARCHAR(128),updated TIMESTAMP,created TIMESTAMP,str VARCHAR(150),int INTEGER)"
)

var (
	unit = schema.Schema{
		Fields: schema.Fields{
			"id": schema.IDField,
			"created": schema.CreatedField,
			"updated": schema.UpdatedField,
			"str": {
				Sortable: true,
				Filterable: true,
				Required: true,
				Validator: &schema.String{
					MaxLen: 150,
				},
			},
			"int": {
				Sortable: true,
				Filterable: true,
				Required: true,
				Validator: &schema.Integer{},
			},
		},
	}
)

func main() {
	s, err := sqlStorage.NewHandler(DB_DRIVER, DB_SOURCE, DB_TABLE)
	if err != nil {
		log.Fatalf("Error connecting database: %s", err)
	}

	err = s.Create(context.TODO(), &unit)
	if err != nil {
		log.Fatalf("Error creating table: %s", err)
	}

	index := resource.NewIndex()
	index.Bind("units", unit, s, resource.Conf{
		AllowedModes: resource.ReadWrite,
	})

	api, err := rest.NewHandler(index)
	if err != nil {
		log.Fatalf("Invalid API configuration: %s", err)
	}

	http.Handle("/", api)

	log.Print("Serving API on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
