package util

import (
	"log"
	"os"
)

func Explode(msg string) {
	log.SetOutput(os.Stdout)

	log.Fatalln(msg)
}
