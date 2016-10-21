package main

import "math/rand"

type name struct {
	first string
	last  string
}

func getFirstNameList() []string {
	firstNames := readLines("names/firstnames.txt")
	dedupeData(&firstNames)
	return firstNames
}

func getLastNameList() []string {
	lastNames := readLines("names/lastnames.txt")
	dedupeData(&lastNames)
	return lastNames
}

func getFirstName() string {
	return app.firstNames[rand.Intn(len(app.firstNames))]
}

func getLastName() string {
	return app.lastNames[rand.Intn(len(app.lastNames))]
}
