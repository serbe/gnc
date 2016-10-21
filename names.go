package main

import "math/rand"

type name struct {
	first string
	last  string
}

func getFirstNameList() {
	firstNames := readLines("names/firstnames.txt")
	dedupeData(&firstNames)
	app.firstNames = firstNames
}

func getLastNameList() {
	lastNames := readLines("names/lastnames.txt")
	dedupeData(&lastNames)
	app.lastNames = lastNames
}

func getFirstName() string {
	return app.firstNames[rand.Intn(len(app.firstNames))]
}

func getLastName() string {
	return app.lastNames[rand.Intn(len(app.lastNames))]
}
