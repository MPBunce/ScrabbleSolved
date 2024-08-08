package main

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "status: available")
	fmt.Fprintf(w, "environment: %s\n", app.config.env)
	fmt.Fprintf(w, "version: %s\n", version)
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	app.render(w, r, "home.page.tmpl")

}

func (app *application) solve(w http.ResponseWriter, r *http.Request) {
	//logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	params := httprouter.ParamsFromContext(r.Context())
	letters := params.ByName("letters")
	if len(letters) > 7 || len(letters) < 1 {
		fmt.Fprintf(w, "Length must be between 1 and 7.")
		return
	}

	lettersLower := strings.ToLower(letters)

	splitWord := strings.Split(lettersLower, "")
	sort.Strings(splitWord)
	sortedWord := strings.Join(splitWord, "")

	if letters == "" {
		app.notFound(w)
		return
	}
	// logger.Printf("logging")
	// if letters == "" {
	// 	app.notFound(w)
	// 	return
	// }

	fmt.Fprintf(w, "Search with letters: %s \n", sortedWord)

}
