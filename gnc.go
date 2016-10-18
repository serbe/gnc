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
	proxyList     []proxy
	currentProxy  int
	words         []string
	goodWordsName string
	badWordsName  string
)

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
	Error  error
	Word   string
}

func getPostString(word string) (string, error) {
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

	// postURL := "https://accounts.google.com/InputValidator?resource=SignUp"
	postURL := "https://173.194.222.84:443/InputValidator?resource=SignUp"

	timeout := time.Duration(10 * time.Second)

	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(urlProxy),
		},
	}

	var data postData

	data.Input01.GmailAddress = word
	data.Input01.Input = "GmailAddress"
	data.Locale = "ru"
	postData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req, _ := http.NewRequest("POST", postURL, bytes.NewBufferString(string(postData)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(string(postData))))
	resp, err := client.Do(req)
	if err != nil {
		if quality > 3 {
			writeLine(host, "proxy/bad_proxy.txt")
			proxyList = append(proxyList[:currentProxy], proxyList[currentProxy+1:]...)
		}
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	jsonResult, err := json.Marshal(string(body))
	if err != nil {
		return "", err
	}

	jsonResultString, err := strconv.Unquote(string(jsonResult))

	return jsonResultString, err
}

func postQuery(word string) postResult {
	var (
		result postResult
	)

	jsonResultString, err := getPostString(word)
	if err != nil {
		result.Error = err
		return result
	}

	err = json.Unmarshal([]byte(jsonResultString), &result)
	result.Word = word
	result.Error = err

	return result
}

func main() {
	numWorkers, err := prepare()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	jobs := make(chan string, len(words))
	results := make(chan postResult, len(words))

	for w := 1; w <= numWorkers; w++ {
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
