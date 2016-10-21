package main

import (
	"flag"
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"time"
)

var (
	app globalValues
)

type globalValues struct {
	firstNames   []string
	lastNames    []string
	words        []string
	proxyList    []proxy
	name         string
	proxyName    string
	length       string
	currentProxy int
	workers      int
	timeout      int
}

func prepare() {
	rand.Seed(time.Now().UnixNano())

	num := runtime.NumCPU()
	runtime.GOMAXPROCS(num)

	wordsName := flag.String("words", "words", "Name of file contains words")
	wordLen := flag.Int("l", 6, "Word length")
	proxyName := flag.String("p", "proxy", "Name of proxy list file")
	numWorkers := flag.Int("w", 15, "Number of workers")
	timeout := flag.Int("t", 15, "Set timeout in second")

	flag.Parse()

	app.name = *wordsName
	app.length = strconv.Itoa(*wordLen)
	app.proxyName = *proxyName
	app.workers = *numWorkers
	app.timeout = *timeout

	// replaceFile("logs.txt")

	getWords()
	getProxyList()
	getFirstNameList()
	getLastNameList()

	fmt.Println("Complete prepare...")
}
