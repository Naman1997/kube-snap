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

	CreateTimedLog(err.Error())
	log.Fatal(err)
}
