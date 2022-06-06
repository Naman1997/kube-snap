package main

import (
	"fmt"
	"os"
	"strings"
)

func getValueOf(key, fallback string) string {
	value, err := os.ReadFile("/etc/secrets/" + key)
	if err != nil {
		return fallback
	}
	return strings.Trim(string(value), "\"")
}

func createFile(path string, data string) {
	path += ".yaml"
	if _, err := os.Stat(path); err == nil {
		os.Remove(path)
		fmt.Println("Re-creating file: ", strings.Replace(path, "/repo/", "", 1))
	}

	f, err := os.Create(path)
	if err != nil {
		panic(err.Error())
	}

	defer f.Close()
	_, err2 := f.WriteString(data)
	if err2 != nil {
		panic(err2.Error())
	}

	fmt.Println("Created: ", path)
}

func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}
