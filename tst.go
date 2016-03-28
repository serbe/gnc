package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type QueryResp struct {
	Addr string
	Time float64
	IP   string
}

var MyIP string

func query(ip string, port string, c chan QueryResp) {
	start_ts := time.Now()
	reIP := regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)
	var timeout = time.Duration(15 * time.Second)
	host := fmt.Sprintf("%s:%s", ip, port)
	url_proxy := &url.URL{Host: host}
	client := &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(url_proxy)},
		Timeout:   timeout}
	resp, err := client.Get("http://myexternalip.com/raw")
	if err != nil {
		c <- QueryResp{Addr: host, Time: float64(-1), IP: ""}
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	time_diff := time.Now().UnixNano() - start_ts.UnixNano()
	if reIP.Match(body) == true {
		c <- QueryResp{Addr: host, Time: float64(time_diff) / 1e9, IP: string(reIP.FindSubmatch(body)[0])}
	} else {
		c <- QueryResp{Addr: host, Time: float64(-1), IP: string(body)}
	}
}

func get_external() string {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func main() {
	MyIP := get_external()
	fmt.Println("my ip is:", MyIP)
	dat, _ := ioutil.ReadFile("ip.txt")
	dats := strings.Split(strings.TrimSuffix(string(dat), "\n"), "\n")

	runtime.GOMAXPROCS(4)

	resp_chan := make(chan QueryResp, 10)

	for _, addr := range dats {
		addrs := strings.SplitN(addr, string(' '), 2)
		ip, port := addrs[0], addrs[1]
		go query(ip, port, resp_chan)
	}

	for _, _ = range dats {
		r := <-resp_chan
		if r.Time > 1e-9 {
			if r.IP != MyIP {
				fmt.Println("addr: " + r.Addr + " time: " + strconv.FormatFloat(r.Time, 'f', -1, 64))
			}
		}
	}
}
