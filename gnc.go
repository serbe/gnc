package main 

import (
    "net/http"
    "encoding/json"
    "fmt"
    "bytes"
    "strconv"
    "io/ioutil"
    "bufio"
    "os"
	"strings"
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

func getResponce(name string) (postResult, error) {
    var ( 
    client http.Client
    data postData
    result postResult
    )
    postURL := "https://accounts.google.com/InputValidator?resource=SignUp"
    data.Input01.GmailAddress = name
    data.Input01.Input = "GmailAddress"
    data.Locale = "ru"
    postJ, err := json.Marshal(data)
    if err != nil {
        return result, err
    }
    req, _ := http.NewRequest("POST", postURL, bytes.NewBufferString(string(postJ)))
	req.Header.Add("Content-Type", "application/json")
    req.Header.Add("Content-Length", strconv.Itoa(len(string(postJ))))
    resp, err := client.Do(req)
    if err != nil {
		return result, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
    jsonResult, err := json.Marshal(string(body)) 
    if err != nil {
		return result, err
	}
    jsonResultString, err := strconv.Unquote(string(jsonResult)) 
    if err != nil {
		return result, err
	}
    
    err = json.Unmarshal([]byte(jsonResultString), &result)
    return result, err
}

func main() {
    lines, err := readLines("words.txt")
    if err != nil {
        panic(err)
    }
    for _, word := range lines {
        if len(word) == 6 {
            responce, err := getResponce(word)
            if err != nil {
                writeLine(word + ", " + fmt.Sprint(err), "err.txt")
            }
            if responce.Input01.Valid == "true" {
                writeLine(word, "valid.txt")
            } else {
                writeLine(word + ", " + strings.Join(responce.Input01.ErrorData, ", "), "novalid.txt")
            }
        }
    }
}