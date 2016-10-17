package main

import "fmt"

func worker(id int, jobs <-chan string, results chan<- postResult) {
	for word := range jobs {
		results <- postQuery(word)
	}
}

func parseResult(r postResult) bool {
	if r.Input01.Valid == "true" {
		writeLine(r.Word, "words/"+App.conf.Name.ValidName+App.conf.Name.Words+"_"+App.strLen+".txt")
		fmt.Println("bingo: ", r.Word)
		return true
	}
	writeLine(fmt.Sprintf("%s %v", r.Word, r.Input01.ErrorData), "words/"+App.conf.Name.NoValidName+App.conf.Name.Words+"_"+App.strLen+".txt")
	return false
}
