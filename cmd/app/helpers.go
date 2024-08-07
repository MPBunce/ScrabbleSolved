package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"sort"
	"strings"
)

func createDic() {
	wordMap := make(map[string][]string)

	file, err := os.Open("sw.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := scanner.Text()
		splitWord := strings.Split(word, "")
		sort.Strings(splitWord)
		wordMap[word] = splitWord
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(wordMap["apple"])
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string) {

	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}
	buf := new(bytes.Buffer)

	err := ts.Execute(buf, nil)

	if err != nil {
		app.serverError(w, err)
		return
	}
	buf.WriteTo(w)

}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.logger.Fatal(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
