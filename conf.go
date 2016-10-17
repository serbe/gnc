package main

import (
	"flag"
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

func prepare() (int, error) {
	var (
		set      map[string]struct{}
		badWords []string
	)

	num := runtime.NumCPU()
	runtime.GOMAXPROCS(num)

	wordsName := flag.String("words", "words", "Name of file contains words")
	wordLen := flag.Int64("l", 6, "Word length")
	proxyName := flag.String("p", "proxy", "Name of proxy list file")
	numWorkers := flag.Int("w", 15, "Number of workers")

	flag.Parse()

	lenStr := strconv.FormatInt(*wordLen, 10)

	goodWordsName = "good_" + *wordsName + "_" + lenStr + ".txt"
	badWordsName = "bad_" + *wordsName + "_" + lenStr + ".txt"

	replaceFile("proxy/" + *proxyName + ".txt")
	existsFile(goodWordsName)
	existsFile(badWordsName)
	existsFile("proxy/bad_proxy.txt")

	lines, err := readLines("words/" + *wordsName + ".txt")
	if err != nil {
		return *numWorkers, err
	}

	goodWordsLines, err := readLines("words/" + goodWordsName)
	if err != nil {
		return *numWorkers, err
	}

	badWordsLines, err := readLines("words/" + badWordsName)
	if err != nil {
		return *numWorkers, err
	}

	fmt.Println("Length of "+badWordsName+" = ", len(badWordsLines))

	for _, line := range badWordsLines {
		array := strings.Split(line, " ")
		if len(array) > 0 {
			badWords = append(badWords, array[0])
		}
	}
	for _, valid := range goodWordsLines {
		if len(valid) == int(*wordLen) {
			badWords = append(badWords, valid)
		}
	}

	replaceFile("words/" + "compact_" + badWordsName)
	writeSlice(badWords, "words/"+"compact_"+badWordsName)
	fmt.Println("Compact "+badWordsName+" = ", len(badWords))

	set = make(map[string]struct{}, len(badWords))
	for _, s := range badWords {
		set[s] = struct{}{}
	}

	for _, word := range lines {
		if len(word) == int(*wordLen) {
			_, ok := set[word]
			if ok == false {
				words = append(words, word)
			}
		}
	}

	fmt.Println("Length of "+*wordsName+".txt = ", len(lines))

	replaceFile("words/" + "compact_" + *wordsName + ".txt")
	writeSlice(words, "words/"+"compact_"+*wordsName+".txt")
	fmt.Println("Compact of "+*wordsName+".txt = ", len(words))

	proxyes, err := getProxyList(*proxyName)
	if err != nil {
		panic(err)
	}
	proxyList = proxyes

	fmt.Println("Num of proxy ", len(proxyList))

	fmt.Println("Complete prepare...")

	return *numWorkers, nil
}
