package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

var (
	proxyList    []proxy
	currentProxy int
	// App application config and more
	App app
)

type app struct {
	config tomlConfig
	strLen string
}

type tomlConfig struct {
	Title     string
	Len       int
	Files     fileNames
	Positions struct {
		Word6 struct {
			letters string
			word    string
		}
		Word7 struct {
			letters string
			word    string
		}
	}
}

type fileNames struct {
	proxyList   string
	goodProxy   string
	badProxy    string
	validName   string
	noValidName string
	words       string
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
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func writeLine(line string, filename string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	fmt.Fprintln(w, line)
	return w.Flush()
}

func existsFile(file string) bool {
	_, err := os.Stat(file)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return createFile(file)
	}
	return true
}

func createFile(file string) bool {
	_, err := os.Create(file)
	if err != nil {
		return false
	}
	return true
}

func (a *app) postQuery(word string) postResult {
	var (
		data    postData
		result  postResult
		timeout = time.Duration(30 * time.Second)
	)

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
			writeLine(host, App.config.Files.badProxy + ".txt")
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

func (a *app) getProxyList() ([]proxy, error) {
	dat, err := ioutil.ReadFile(App.config.Files.proxyList + ".txt")
	dats := strings.Split(strings.TrimSuffix(string(dat), "\n"), "\n")
	var proxyList []proxy
	for _, host := range dats {
		var tmpproxy proxy
		tmpproxy.host = host
		proxyList = append(proxyList, tmpproxy)
	}
	return proxyList, err
}

func main() {
	var (
		twoLetter string
		valid     int
	)

	if _, err := toml.DecodeFile("example.toml", &App.config); err != nil {
		panic(err)
	}

	App.strLen = fmt.Sprintf("%d", App.config.Len)

	os.Remove(App.config.Files.goodProxy + ".txt")
	createFile(App.config.Files.goodProxy + ".txt")
	existsFile(App.config.Files.validName + "valid" + App.strLen + ".txt")
	existsFile(App.config.Files.noValidName + "novalid" + App.strLen + ".txt")
	existsFile(App.config.Files.badProxy + ".txt")

	runtime.GOMAXPROCS(4)

	lines, err := readLines(App.config.Files.badProxy + ".txt")
	if err != nil {
		panic(err)
	}

	proxyes, err := App.getProxyList()
	if err != nil {
		panic(err)
	}
	proxyList = proxyes

	for _, word := range lines {
		if len(word) == App.config.Len {
			var (
				r postResult
			)
			r = App.postQuery(word)
			letters := word[0:2]
			if twoLetter != letters {
				fmt.Println(letters)
				twoLetter = letters
			}
			if r.Input01.Valid == "true" {
				valid++
				writeLine(word, App.config.Files.validName + "valid" + App.strLen + ".txt")
				fmt.Println("bingo: ", word)
				if valid == 10 {
					panic(err)
				}
			} else {
				valid = 0
				writeLine(fmt.Sprintf("%s %v", word, r.Input01.ErrorData), App.config.Files.noValidName + "novalid" + App.strLen + ".txt")
			}
		}
	}
	for _, i := range proxyList {
		writeLine(i.host, App.config.Files.goodProxy + ".txt")
	}
}
