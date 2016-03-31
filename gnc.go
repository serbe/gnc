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
)

var (
	proxyList    []proxy
	currentProxy int
	// App application config and more
	App app
)

type app struct {
	conf config
	strLen string
}

type config struct {
	Len  int `json:"len"`
	Name struct {
		BadProxy    string `json:"bad_proxy"`
		GoodProxy   string `json:"good_proxy"`
		NoValidName string `json:"no_valid_name"`
		ProxyList   string `json:"proxy_list"`
		ValidName   string `json:"valid_name"`
		Words       string `json:"words"`
	} `json:"name"`
	Position struct {
		Word6 struct {
			Letters string `json:"letters"`
			Word    string `json:"word"`
		} `json:"word6"`
		Word7 struct {
			Letters string `json:"letters"`
			Word    string `json:"word"`
		} `json:"word7"`
	} `json:"position"`
	Title string `json:"title"`
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

func getConfig() (config, error) {
	c := config{}
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		return c, err
	}
	err = json.Unmarshal(file, &c)
	return c, err
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
			writeLine(host, App.conf.Name.BadProxy+".txt")
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
	dat, err := ioutil.ReadFile(App.conf.Name.ProxyList + ".txt")
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
	conf, err := getConfig()
    if err != nil {
		fmt.Println(err)
		panic(err)
	}

    App.conf = conf
	App.strLen = fmt.Sprintf("%d", App.conf.Len)

	os.Remove(App.conf.Name.GoodProxy + ".txt")
	createFile(App.conf.Name.GoodProxy + ".txt")
	existsFile(App.conf.Name.ValidName + App.strLen + ".txt")
	existsFile(App.conf.Name.NoValidName + App.strLen + ".txt")
	existsFile(App.conf.Name.BadProxy + ".txt")

	runtime.GOMAXPROCS(4)

	lines, err := readLines(App.conf.Name.Words + ".txt")
	if err != nil {
		panic(err)
	}

	proxyes, err := App.getProxyList()
	if err != nil {
		panic(err)
	}
	proxyList = proxyes

	for _, word := range lines {
		if len(word) == App.conf.Len {
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
				writeLine(word, App.conf.Name.ValidName+"valid"+App.strLen+".txt")
				fmt.Println("bingo: ", word)
				if valid == 10 {
					panic(err)
				}
			} else {
				valid = 0
				writeLine(fmt.Sprintf("%s %v", word, r.Input01.ErrorData), App.conf.Name.NoValidName+"novalid"+App.strLen+".txt")
			}
		}
	}
	for _, i := range proxyList {
		writeLine(i.host, App.conf.Name.GoodProxy+".txt")
	}
}
