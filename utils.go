package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var (
	reProxy = regexp.MustCompile(`.*?(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d{2,5}).*?`)
)

func removeDuplicates(dats []string) []string {
	var newDats []string
	for _, line := range dats {
		if reProxy.MatchString(line) {
			newDats = append(newDats, reProxy.FindString(line))
		}
	}
	parseProxy(&newDats)
	return newDats
}

func parseProxy(xs *[]string) {
	found := make(map[string]bool)
	j := 0
	for i, x := range *xs {
		if !found[x] {
			found[x] = true
			(*xs)[j] = (*xs)[i]
			j++
		}
	}
	*xs = (*xs)[:j]
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

func writeSlice(slice []string, filename string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range slice {
		fmt.Fprintln(w, line)
	}
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

func replaceFile(file string) bool {
	err := os.Remove(file)
	if err != nil {
		return false
	}
	return existsFile(file)
}

func getProxyList(name string) ([]proxy, error) {
	dat, err := ioutil.ReadFile("proxy/" + name + ".txt")
	dats := strings.Split(strings.TrimSuffix(string(dat), "\n"), "\n")
	dats = removeDuplicates(dats)
	var proxyList []proxy
	for _, host := range dats {
		var tmpproxy proxy
		tmpproxy.host = host
		proxyList = append(proxyList, tmpproxy)
	}
	return proxyList, err
}
