package main

import (
	"flag"
	"fmt"
	"runtime"
	"strconv"
)

func prepare() {
	num := runtime.NumCPU()
	runtime.GOMAXPROCS(num)

	wordsName := flag.String("words", "words", "Name of file contains words")
	wordLen := flag.Int("l", 6, "Word length")
	proxyName := flag.String("p", "proxy", "Name of proxy list file")
	numWorkers := flag.Int("w", 15, "Number of workers")

	flag.Parse()

	app.name = *wordsName
	app.length = strconv.Itoa(*wordLen)
	app.proxyName = *proxyName
	app.workers = *numWorkers

	replaceFile("logs.txt")

	app.words = getWords()
	app.proxyList = getProxyList()
	app.firstNames = getFirstNameList()
	app.lastNames = getLastNameList()

	fmt.Println("Complete prepare...")
}
