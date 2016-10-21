package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var (
	reProxy = regexp.MustCompile(`.*?(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d{2,5}).*?`)
)

func createProxyFiles() {
	replaceFile("proxy/good_" + app.proxyName + ".txt")
}

func cleanProxyList(stringList []string) []string {
	var data []string
	for _, line := range stringList {
		if reProxy.MatchString(line) {
			data = append(data, reProxy.FindString(line))
		}
	}
	dedupeData(&data)
	return data
}

func getProxyList() []string {
	createProxyFiles()
	data, err := ioutil.ReadFile("proxy/" + app.proxyName + ".txt")
	if err != nil {
		panic(err)
	}
	proxyList := strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
	proxyList = cleanProxyList(proxyList)
	fmt.Println("Num of proxy ", len(proxyList))
	return proxyList
}

func getProxy() string {
	lenProxy := len(app.proxyList)
	if lenProxy <= 1 {
		fmt.Println("end proxy list")
		os.Exit(1)
	}
	app.currentProxy++
	if app.currentProxy >= lenProxy {
		app.currentProxy = 0
	}
	return app.proxyList[app.currentProxy]
}
