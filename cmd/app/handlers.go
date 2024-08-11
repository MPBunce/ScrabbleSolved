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
	data := templateData{
		SearchWord: "", // Initialize with an empty string
		Matches:    nil,
		OneAway:    nil,
		TwoAway:    nil,
	}
	app.render(w, r, "home.page.tmpl", &data)

}

func (app *application) solve(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())
	letters := params.ByName("searchWord")
	fmt.Println(letters)
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

	//fmt.Fprintf(w, "Search with letters: %s \n", sortedWord)

	//Can Make
	var canMakeArray = []scrabbleWords{}
	res := getCurrent(sortedWord)
	for _, v := range res {
		val, ok := app.scrabbleWords[v]
		// If the key exists
		if ok {
			for _, v := range val {
				canMakeArray = append(canMakeArray, v)
			}
		}
	}

	//print to screen
	// for _, v := range canMakeArray {
	// 	fmt.Fprintf(w, "%s %d\n", v.Word, v.Score)
	// }

	//One Away
	//fmt.Fprintf(w, "OneAway\n\n ")
	var oneAwayArray = []scrabbleWords{}
	resOneAway := getOneAway(sortedWord)
	for _, v := range resOneAway {
		val, ok := app.scrabbleWords[v]
		// If the key exists
		if ok {
			for _, v := range val {
				if !containsScrabbleWord(canMakeArray, v.Word) {
					oneAwayArray = append(oneAwayArray, v)
				}
			}
		}
	}
	//print to screen
	// for _, v := range oneAwayArray {
	// 	fmt.Fprintf(w, "%s %d\n", v.Word, v.Score)
	// }

	//Two Away

	//fmt.Fprintf(w, "TwoAway\n\n ")
	var twoAwayArray = []scrabbleWords{}
	resTwoAway := getTwoAway(sortedWord)
	for _, v := range resTwoAway {
		val, ok := app.scrabbleWords[v]
		// If the key exists
		if ok {
			for _, v := range val {
				if !containsScrabbleWord(canMakeArray, v.Word) {
					if !containsScrabbleWord(oneAwayArray, v.Word) {
						twoAwayArray = append(twoAwayArray, v)
					}
				}
			}
		}
	}

	//print to screen
	// for _, v := range twoAwayArray {
	// 	fmt.Fprintf(w, "%s %d\n", v.Word, v.Score)
	// }

	app.render(w, r, "show.page.tmpl", &templateData{
		SearchWord: letters,
		Matches:    canMakeArray,
		OneAway:    oneAwayArray,
		TwoAway:    twoAwayArray,
	})

	// dummyData := &templateData{
	// 	SearchWord: "abc",
	// 	Matches: []scrabbleWords{
	// 		{Word: "cab", Score: 8},
	// 		{Word: "bac", Score: 8},
	// 	},
	// 	OneAway: []scrabbleWords{
	// 		{Word: "ab", Score: 5},
	// 		{Word: "bc", Score: 6},
	// 	},
	// 	TwoAway: []scrabbleWords{
	// 		{Word: "a", Score: 1},
	// 		{Word: "b", Score: 3},
	// 		{Word: "c", Score: 3},
	// 	},
	// }

	//app.render(w, r, "show.page.tmpl", dummyData)
}
