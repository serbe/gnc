package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
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

func getPostBody(word string) ([]byte, error) {
	proxy := getProxy()

	postURL := "https://accounts.google.com/InputValidator?resource=SignUp"
	// postURL := "https://64.233.164.84:443/InputValidator?resource=SignUp"

	client := &http.Client{
		Timeout: time.Duration(app.timeout) * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(&url.URL{
				Host: proxy.host,
			}),
		},
	}

	postString, _ := getPostString(word)

	req, _ := http.NewRequest("POST", postURL, bytes.NewBufferString(postString))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(postString)))
	resp, err := client.Do(req)
	if err != nil {
		removeProxy(&app.proxyList, proxy)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		removeProxy(&app.proxyList, proxy)
		return nil, err
	}
	return body, nil
}

func getResponse(word string) string {
	var (
		successJSON bool
		jsonBytes   []byte
	)
	for !successJSON {
		var (
			successGetBody bool
			body           []byte
			err            error
		)
		for !successGetBody {
			// fmt.Println("tryGetBody", word)
			body, err = getPostBody(word)
			if err == nil {
				successGetBody = true
			} else {
				// fmt.Println("error tryGetBody", err)
			}
		}
		// fmt.Println("jsonBytes", word)
		jsonBytes, err = json.Marshal(string(body))
		if err == nil {
			successJSON = true
		} else {
			// fmt.Println("error jsonBytes", err)
		}
	}

	jsonString, _ := strconv.Unquote(string(jsonBytes))

	return jsonString
}

func postQuery(word string) postResult {
	var (
		result postResult
	)

	response := getResponse(word)

	err := json.Unmarshal([]byte(response), &result)
	result.Word = word
	result.Error = err

	return result
}
