package main

import (
	"log"
	"os"
)

const MEME_FOLDER = "./static/"
const JACHOB_ROLE = "1295953082827407392"

var infolog = log.New(os.Stdout, "INFO  : ", log.LstdFlags|log.Lshortfile)
var errorlog = log.New(os.Stdout, "ERROR : ", log.LstdFlags|log.Lshortfile)

func checkErr(err error) {
	if err != nil {
		errorlog.Fatal(err)
	}
}
