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

	WordMap := make(map[string]scrabbleLetter)
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

		WordMap[splitWord[0]] = letter

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return WordMap
}

func createWordsDic(filePath string, scrabbleLetters map[string]scrabbleLetter) map[string][]scrabbleWords {

	WordMap := make(map[string][]scrabbleWords)

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		Word := scanner.Text()
		splitWord := strings.Split(Word, "")
		sort.Strings(splitWord)
		sortedWord := strings.Join(splitWord, "")

		calcScore := 0
		for _, char := range sortedWord {
			calcScore += scrabbleLetters[string(char)].value
		}

		data := scrabbleWords{
			Word:  Word,
			Score: calcScore,
		}

		WordMap[sortedWord] = append(WordMap[sortedWord], data)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return WordMap
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {

	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}
	buf := new(bytes.Buffer)

	// Pass td (templateData) instead of nil
	err := ts.Execute(buf, td)

	if err != nil {
		app.serverError(w, err)
		return
	}
	buf.WriteTo(w)
}

func getOneAway(input string) (output []string) {
	letters := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

	freq := countFrequencies(input)
	results := []string{}

	for _, v := range letters {
		generateCombinations(input, v, freq, &results)
	}

	for _, value := range results {

		splitWord := strings.Split(value, "")
		sort.Strings(splitWord)
		sortedWord := strings.Join(splitWord, "")

		if !contains(output, sortedWord) {
			output = append(output, sortedWord)
		}

	}

	return
}

func getTwoAway(input string) (output []string) {
	letters := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

	freq := countFrequencies(input)
	results := []string{}

	for _, v := range letters {
		for _, val := range letters {
			str := []string{v, val}
			temp := strings.Join(str, "")
			//fmt.Println(temp)
			generateCombinations(input, temp, freq, &results)
		}

	}

	for _, value := range results {

		splitWord := strings.Split(value, "")
		sort.Strings(splitWord)
		sortedWord := strings.Join(splitWord, "")

		if !contains(output, sortedWord) {
			output = append(output, sortedWord)
		}

	}

	return
}

func getCurrent(input string) (output []string) {

	freq := countFrequencies(input)
	results := []string{}
	generateCombinations(input, "", freq, &results)

	for _, value := range results {

		splitWord := strings.Split(value, "")
		sort.Strings(splitWord)
		sortedWord := strings.Join(splitWord, "")

		if !contains(output, sortedWord) {
			output = append(output, sortedWord)
		}

	}

	return
}

func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func containsScrabbleWord(slice []scrabbleWords, value string) bool {
	for _, v := range slice {
		//fmt.Printf("%s compare %s\n", v.Word, value)
		if v.Word == value {
			return true
		}
	}
	return false
}

func countFrequencies(input string) map[rune]int {
	freq := make(map[rune]int)
	for _, char := range input {
		freq[char]++
	}
	return freq
}

func generateCombinations(input string, prefix string, freq map[rune]int, results *[]string) {

	if len(prefix) > 0 {
		*results = append(*results, prefix)
	}

	for char, count := range freq {
		if count > 0 {
			newFreq := make(map[rune]int)
			for k, v := range freq {
				newFreq[k] = v
			}
			newFreq[char]--
			newPrefix := prefix + string(char)
			generateCombinations(input, newPrefix, newFreq, results)
		}
	}

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
