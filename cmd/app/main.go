package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {

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
		wordMap[word] = splitWord
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(wordMap["apple"])
}
