package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

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
	var data postData

	data.Input01.FirstName = getFirstName()
	data.Input01.LastName = getLastName()
	data.Input01.GmailAddress = word
	data.Input01.Input = "GmailAddress"
	data.Locale = "ru"
	pbyte, err := json.Marshal(data)
	return string(pbyte), err
}

func getPost(word string) (string, error) {
	proxy := getProxy()

	postURL := "https://accounts.google.com/InputValidator?resource=SignUp"
	// postURL := "https://64.233.164.84:443/InputValidator?resource=SignUp"

	timeout := time.Duration(10 * time.Second)

	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(&url.URL{
				Host: proxy,
			}),
		},
	}

	postString, err := getPostString(word)

	req, _ := http.NewRequest("POST", postURL, bytes.NewBufferString(postString))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(postString)))
	resp, err := client.Do(req)
	if err != nil {
		// if quality > 3 {
		// 	writeLine(host, "proxy/bad_proxy.txt")
		// 	proxyList = append(proxyList[:currentProxy], proxyList[currentProxy+1:]...)
		// }
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	writeLine(proxy+" : "+word+"\n", "log.txt")
	writeLine(postString+"\n", "log.txt")
	writeLine(string(body)+"\n", "log.txt")

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

	postString, err := getPost(word)
	if err != nil {
		log.Println(err)
		result.Error = err
		return result
	}

	err = json.Unmarshal([]byte(postString), &result)
	result.Word = word
	result.Error = err

	return result
}
