package utilities

import (
	"log"
)

func CheckIfError(err error, message string) {
	if err == nil {
		return
	}

	if message != "" {
		CreateTimedLog("[ERROR]", message)
	}

	log.Fatal(err)
}
