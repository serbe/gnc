package main

var (
	app globalValues
)

type globalValues struct {
	firstNames   []string
	lastNames    []string
	words        []string
	proxyList    []string
	name         string
	proxyName    string
	length       string
	currentProxy int
	workers      int
}

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
}
