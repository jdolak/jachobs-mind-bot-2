package main

import (
	"log"
)

const MEME_FOLDER = "./static/"
const JACHOB_ROLE = "1295953082827407392"

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
