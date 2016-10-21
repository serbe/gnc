package main

import (
	"fmt"
	"strconv"
	"strings"
)

func createWordFiles() {
	existsFile("words/good_" + app.name + "_" + app.length + ".txt")
	existsFile("words/bad_" + app.name + "_" + app.length + ".txt")
	existsFile("words/hint_" + app.name + "_" + app.length + ".txt")
	replaceFile("words/" + "compact_bad_" + app.name + "_" + app.length + ".txt")
	replaceFile("words/" + "compact_" + app.name + ".txt")
}

func getWords() {
	var (
		badWords []string
		words    []string
		set      map[string]struct{}
	)

	length, err := strconv.Atoi(app.length)
	if err != nil {
		panic(err)
	}

	createWordFiles()

	lines := readLines("words/" + app.name + ".txt")
	goodWordsLines := readLines("words/good_" + app.name + "_" + app.length + ".txt")
	badWordsLines := readLines("words/bad_" + app.name + "_" + app.length + ".txt")

	fmt.Println("Length of bad_"+app.name+"_"+app.length+".txt = ", len(badWordsLines))

	for _, line := range badWordsLines {
		array := strings.Split(line, " ")
		if len(array) > 0 {
			badWords = append(badWords, array[0])
		}
	}
	for _, valid := range goodWordsLines {
		if len(valid) == length {
			badWords = append(badWords, valid)
		}
	}

	writeSlice(badWords, "words/"+"compact_bad_"+app.name+"_"+app.length+".txt")
	fmt.Println("Compact bad_"+app.name+"_"+app.length+".txt = ", len(badWords))

	set = make(map[string]struct{}, len(badWords))
	for _, s := range badWords {
		set[s] = struct{}{}
	}

	for _, word := range lines {
		if len(word) == length {
			_, ok := set[word]
			if ok == false {
				words = append(words, word)
			}
		}
	}

	fmt.Println("Length of "+app.name+".txt = ", len(lines))
	writeSlice(words, "words/"+"compact_"+app.name+".txt")

	app.words = words

	fmt.Println("Compact of "+app.name+".txt = ", len(app.words))
}
