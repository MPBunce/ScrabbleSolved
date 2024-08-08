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
	"strconv"
	"strings"
)

func createLetterDic(filePath string) map[string]scrabbleLetter {

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	wordMap := make(map[string]scrabbleLetter)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		line := scanner.Text()
		splitWord := strings.Split(line, "-")
		if len(splitWord) != 3 {
			log.Fatal("Invalid line format")
		}

		value, err := strconv.Atoi(splitWord[1])
		if err != nil {
			log.Fatal(err)
		}

		count, err := strconv.Atoi(splitWord[2])
		if err != nil {
			log.Fatal(err)
		}

		letter := scrabbleLetter{
			value: value,
			count: count,
		}

		wordMap[splitWord[0]] = letter

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return wordMap
}

func createWordsDic(filePath string, scrabbleLetters map[string]scrabbleLetter) map[string][]scrabbleWords {

	wordMap := make(map[string][]scrabbleWords)

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := scanner.Text()
		splitWord := strings.Split(word, "")
		sort.Strings(splitWord)
		sortedWord := strings.Join(splitWord, "")

		calcScore := 0
		for _, char := range sortedWord {
			calcScore += scrabbleLetters[string(char)].value
		}

		data := scrabbleWords{
			word:  word,
			score: calcScore,
		}

		wordMap[sortedWord] = append(wordMap[sortedWord], data)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return wordMap
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
