package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

var (
	reProxy = regexp.MustCompile(`.*?(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d{2,5}).*?`)
)

type proxy struct {
	host string
	good bool
}

func createProxyFiles() {
	replaceFile("proxy/good_" + app.proxyName + ".txt")
	replaceFile("proxy/bad_" + app.proxyName + ".txt")
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

func getProxyList() {
	var list []proxy

	createProxyFiles()
	data, err := ioutil.ReadFile("proxy/" + app.proxyName + ".txt")
	if err != nil {
		panic(err)
	}
	proxyList := strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
	proxyList = cleanProxyList(proxyList)
	fmt.Println("Num of proxy ", len(proxyList))

	for _, l := range proxyList {
		var pr proxy
		pr.host = l
		pr.good = true
		list = append(list, pr)
	}

	app.proxyList = list
}

func getProxy() proxy {
	lenProxy := len(app.proxyList)
	var findGoodProxy bool
	for !findGoodProxy {
		app.currentProxy++
		if app.currentProxy >= lenProxy {
			app.currentProxy = 0
		}
		if app.proxyList[app.currentProxy].good == true {
			findGoodProxy = true
		}
	}
	return app.proxyList[app.currentProxy]
}

func deleteProxy(proxy proxy) {
	for i, p := range app.proxyList {
		if p.host == proxy.host {
			app.proxyList[i].good = false
		}
	}
	writeLine(proxy.host, "proxy/bad_"+app.proxyName+".txt")
}

func saveProxyList() {
	for _, p := range app.proxyList {
		if p.good {
			writeLine(p.host, "proxy/good_"+app.proxyName+".txt")
		}
	}
}
