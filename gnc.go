package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

var (
	proxyList    []proxy
	currentProxy int
	// App application config and more
	App          app
	set          map[string]struct{}
	numWords     int
	curNum       int
	proc         int
	noValidWords []string
	words        []string
)

type app struct {
	conf   config
	strLen string
	work   int
}

type proxy struct {
	host    string
	quality int
}

type postData struct {
	Input01 struct {
		FirstName    string `json:"FirstName"`
		GmailAddress string `json:"GmailAddress"`
		Input        string `json:"Input"`
		LastName     string `json:"LastName"`
	} `json:"input01"`
	Locale string `json:"Locale"`
}

type postResult struct {
	Input01 struct {
		ErrorData    []string `json:"ErrorData"`
		ErrorMessage string   `json:"ErrorMessage"`
		Errors       struct {
			GmailAddress string `json:"GmailAddress"`
		} `json:"Errors"`
		Valid string `json:"Valid"`
	} `json:"input01"`
	Locale string `json:"Locale"`
	Proxy  string
	Error  error
	Word   string
}

func postQuery(word string) postResult {
	var (
		data    postData
		result  postResult
		timeout = time.Duration(10 * time.Second)
	)

	result.Word = word
	lenProxy := len(proxyList)
	if lenProxy <= 1 {
		fmt.Println("end proxy list")
		os.Exit(1)
	}
	currentProxy++
	if currentProxy >= lenProxy {
		currentProxy = 0
	}
	host := proxyList[currentProxy].host
	quality := proxyList[currentProxy].quality

	urlProxy := &url.URL{Host: host}
	result.Proxy = host

	postURL := "https://accounts.google.com/InputValidator?resource=SignUp"
	// postURL := "https://64.233.162.139/InputValidator?resource=SignUp"
	data.Input01.GmailAddress = word
	data.Input01.Input = "GmailAddress"
	data.Locale = "ru"
	postJ, err := json.Marshal(data)
	if err != nil {
		result.Error = err
		return result
	}

	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(urlProxy),
		},
	}

	req, _ := http.NewRequest("POST", postURL, bytes.NewBufferString(string(postJ)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(string(postJ))))
	resp, err := client.Do(req)
	if err != nil {
		if quality > 3 {
			writeLine(host, "words/"+App.conf.Name.BadProxy+".txt")
			proxyList = append(proxyList[:currentProxy], proxyList[currentProxy+1:]...)

		}
		result.Error = err
		return result
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	jsonResult, err := json.Marshal(string(body))
	if err != nil {
		result.Error = err
		return result
	}
	jsonResultString, err := strconv.Unquote(string(jsonResult))
	if err != nil {
		result.Error = err
		return result
	}

	err = json.Unmarshal([]byte(jsonResultString), &result)
	result.Error = err
	return result
}

func main() {
	err := prepare()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	numWords = len(words)
	curNum = 0
	proc = 0

	jobs := make(chan string, len(words))
	results := make(chan postResult, len(words))

	for w := 1; w <= App.conf.Workers; w++ {
		go worker(w, jobs, results)
	}

	for _, word := range words {
		jobs <- word
	}

	for _ = range words {
		r := <-results
		parseResult(r)
	}
}
