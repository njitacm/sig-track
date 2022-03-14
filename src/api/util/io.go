package util

import (
	"log"
	"os"
)

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func FileWrite(filename string, obj []byte) {
	f, err := os.Create(filename)
	Check(err)

	_, err = f.Write(obj)
	Check(err)
}
