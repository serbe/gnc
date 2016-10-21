package main

func main() {
	prepare()

	jobs := make(chan string, len(app.words))
	results := make(chan postResult, len(app.words))

	for w := 1; w <= app.workers; w++ {
		go worker(w, jobs, results)
	}

	for _, word := range app.words {
		jobs <- word
	}

	for _ = range app.words {
		r := <-results
		parseResult(r)
	}

	saveProxyList()
}
