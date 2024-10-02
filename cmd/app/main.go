package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type scrabbleWords struct {
	Word  string
	Score int
}

type scrabbleLetter struct {
	value int
	count int
}

type application struct {
	config         config
	logger         *log.Logger
	scrabbleWords  map[string][]scrabbleWords
	scrabbleLetter map[string]scrabbleLetter
	templateCache  map[string]*template.Template
}

func main() {

	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		logger.Fatal(err)
	}

	scrabbleLetters := createLetterDic("./ui/static/letterData.txt")
	if err != nil {
		logger.Fatal(err)
	}

	//letterA := scrabbleLetters["a"]
	//logger.Printf("Letter: %s, Value: %d, Count: %d\n", "a", letterA.value, letterA.count)

	scrabbleWords := createWordsDic("./ui/static/wordData.txt", scrabbleLetters)
	if err != nil {
		logger.Fatal(err)
	}
	//logger.Println(scrabbleWords["fmor"])
	//Count to confirm
	// count := 0
	// for letter, words := range scrabbleWords {
	// 	fmt.Println("Letter:", letter)
	// 	for range words {
	// 		count += 1
	// 	}
	// }
	// logger.Printf("count: %d\n", count)

	app := &application{
		config:         cfg,
		logger:         logger,
		scrabbleLetter: scrabbleLetters,
		scrabbleWords:  scrabbleWords,
		templateCache:  templateCache,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)

}
