package main

import (
	"bufio"
	"fmt"
	"os"
)

func dedupeData(xs *[]string) {
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

func removeData(xs *[]string, data string) {
	j := 0
	for i, x := range *xs {
		if x != data {
			(*xs)[j] = (*xs)[i]
			j++
		}
	}
	*xs = (*xs)[:j]
}

func readLines(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if len(scanner.Text()) > 0 {
			lines = append(lines, scanner.Text())
		}
	}
	return lines
}

func writeLine(line string, filename string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	fmt.Fprintln(w, line)
	return w.Flush()
}

func writeSlice(slice []string, filename string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
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
	_ = os.Remove(file)
	return existsFile(file)
}

func parseResult(r postResult) bool {
	if r.Input01.Valid == "true" {
		writeLine(r.Word, "words/good_"+app.name+"_"+app.length+".txt")
		fmt.Println("bingo: ", r.Word)
		return true
	}
	writeLine(fmt.Sprintf("%s", r.Word), "words/bad_"+app.name+"_"+app.length+".txt")
	if len(r.Input01.ErrorData) > 0 {
		for _, w := range r.Input01.ErrorData {
			writeLine(fmt.Sprintf("%s", w), "words/hint_"+app.name+"_"+app.length+".txt")
		}
	}
	return false
}
