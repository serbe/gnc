package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"
)

type config struct {
	Len     int `json:"len"`
	Workers int `json:"workers"`
	Name    struct {
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

func getConfig() (config, error) {
	c := config{}
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		return c, err
	}
	err = json.Unmarshal(file, &c)
	return c, err
}

func prepare() error {
	num := runtime.NumCPU()
	runtime.GOMAXPROCS(num)

	conf, err := getConfig()
	if err != nil {
		fmt.Println(err)
		return err
	}

	App.conf = conf
	App.strLen = fmt.Sprintf("%d", App.conf.Len)

	replaceFile("proxy/" + App.conf.Name.GoodProxy + ".txt")
	existsFile("words/" + App.conf.Name.ValidName + App.conf.Name.Words + "_" + App.strLen + ".txt")
	existsFile("words/" + App.conf.Name.NoValidName + App.conf.Name.Words + "_" + App.strLen + ".txt")
	existsFile("proxy/" + App.conf.Name.BadProxy + ".txt")

	lines, err := readLines("words/" + App.conf.Name.Words + ".txt")
	if err != nil {
		return err
	}

	valids, err := readLines("words/" + App.conf.Name.ValidName + App.conf.Name.Words + "_" + App.strLen + ".txt")
	if err != nil {
		return err
	}

	noValidWordsLines, err := readLines("words/" + App.conf.Name.NoValidName + App.conf.Name.Words + "_" + App.strLen + ".txt")
	if err != nil {
		return err
	}

	fmt.Println("Length of "+App.conf.Name.NoValidName+App.conf.Name.Words+"_"+App.strLen+".txt = ", len(noValidWordsLines))

	for _, line := range noValidWordsLines {
		array := strings.Split(line, " ")
		if len(array) > 0 {
			noValidWords = append(noValidWords, array[0])
		}
	}
	for _, valid := range valids {
		if len(valid) == App.conf.Len {
			noValidWords = append(noValidWords, valid)
		}
	}

	replaceFile("words/" + "compact" + App.conf.Name.NoValidName + App.conf.Name.Words + "_" + App.strLen + ".txt")
	writeSlice(noValidWords, "words/"+"compact"+App.conf.Name.NoValidName+App.conf.Name.Words+"_"+App.strLen+".txt")
	fmt.Println("Compact "+App.conf.Name.NoValidName+App.conf.Name.Words+"_"+App.strLen+".txt = ", len(noValidWords))

	set = make(map[string]struct{}, len(noValidWords))
	for _, s := range noValidWords {
		set[s] = struct{}{}
	}

	for _, word := range lines {
		if len(word) == App.conf.Len {
			_, ok := set[word]
			if ok == false {
				words = append(words, word)
			}
		}
	}

	fmt.Println("Length of "+App.conf.Name.Words+".txt = ", len(lines))

	replaceFile("words/" + "compact" + App.conf.Name.Words + ".txt")
	writeSlice(words, "words/"+"compact"+App.conf.Name.Words+".txt")
	fmt.Println("Compact of "+App.conf.Name.Words+".txt = ", len(words))

	proxyes, err := getProxyList()
	if err != nil {
		panic(err)
	}
	proxyList = proxyes

	fmt.Println("Num of proxy ", len(proxyList))

	fmt.Println("Complete prepare...")

	return nil
}
