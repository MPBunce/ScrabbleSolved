package main

import (
	"fmt"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "status: available")
	fmt.Fprintf(w, "environment: %s\n", app.config.env)
	fmt.Fprintf(w, "version: %s\n", version)
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	app.render(w, r, "home.page.tmpl")

}

func (app *application) search(w http.ResponseWriter, r *http.Request) {

	app.render(w, r, "home.page.tmpl")

}
