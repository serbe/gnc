package main

func worker(id int, jobs <-chan string, results chan<- postResult) {
	for word := range jobs {
		results <- postQuery(word)
	}
}
