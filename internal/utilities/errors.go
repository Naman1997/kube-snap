package utilities

import (
	"fmt"
	"log"
)

func CheckIfError(err error, message string) {
	if err == nil {
		return
	}

	if message != "" {
		fmt.Println("[ERROR] ", message)
	}

	log.Fatal(err)
}
