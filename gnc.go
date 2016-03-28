package main 

import (
    "net/http"
    "encoding/json"
    "fmt"
    "bytes"
    "strconv"
    "io/ioutil"
)

type postData struct {
	Input01 struct {
		FirstName    string `json:"FirstName"`
		GmailAddress string `json:"GmailAddress"`
		Input        string `json:"Input"`
		LastName     string `json:"LastName"`
	} `json:"input01"`
    Locale  string `json:"Locale"`
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
   	Locale  string `json:"Locale"`
}

func main() {
    var client http.Client
    var data postData
    postURL := "https://accounts.google.com/InputValidator?resource=SignUp"
    data.Input01.GmailAddress = "oldhurd"
    data.Input01.Input = "GmailAddress"
    data.Locale = "ru"
    postJson, err := json.Marshal(data)
    if err != nil {
        panic(err)
    }
    req, _ := http.NewRequest("POST", postURL, bytes.NewBufferString(string(postJson)))
	req.Header.Add("Content-Type", "application/json")
    req.Header.Add("Content-Length", strconv.Itoa(len(string(postJson))))
    resp, err := client.Do(req)
    if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
    jsonResult, err := json.Marshal(string(body)) 
    if err != nil {
		panic(err)
	}
    jsonResultString, err := strconv.Unquote(string(jsonResult)) 
    if err != nil {
		panic(err)
	}
    var result postResult
    err = json.Unmarshal([]byte(jsonResultString), &result)
    fmt.Println(result.Input01.ErrorData, result.Input01.Valid)
}