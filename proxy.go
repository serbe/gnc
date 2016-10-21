package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"time"
)

var (
	reProxy = regexp.MustCompile(`.*?(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d{2,5}).*?`)
)

type proxy struct {
	host string
	time time.Time
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
		pr.time = time.Now()
		list = append(list, pr)
	}

	app.proxyList = list
}

func getProxy() proxy {
	mutex.Lock()
	lenProxy := len(app.proxyList)
	var findGoodProxy bool
	for !findGoodProxy {
		app.currentProxy++
		if app.currentProxy >= lenProxy {
			app.currentProxy = 0
		}
		t := time.Now()
		if t.Sub(app.proxyList[app.currentProxy].time) > time.Duration(10*time.Second) {
			findGoodProxy = true
		}
	}
	app.proxyList[app.currentProxy].time = time.Now()
	mutex.Unlock()
	return app.proxyList[app.currentProxy]
}

func saveProxyList() {
	for _, p := range app.proxyList {
		writeLine(p.host, "proxy/good_"+app.proxyName+".txt")
	}
}

func removeProxy(proxyList *[]proxy, proxy proxy) {
	mutex.Lock()
	j := 0
	for i, x := range *proxyList {
		if x.host != proxy.host {
			(*proxyList)[j] = (*proxyList)[i]
			j++
		}
	}
	*proxyList = (*proxyList)[:j]
	mutex.Unlock()
	writeLine(proxy.host, "proxy/bad_"+app.proxyName+".txt")
}
