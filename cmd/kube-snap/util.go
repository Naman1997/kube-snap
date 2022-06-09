package main

import (
	"fmt"
	"os"
	"strings"
)

func getValueOf(key string) string {
	value, _ := os.ReadFile("/etc/secrets/" + key)
	updatedValue := strings.Trim(string(value), "\"")
	return strings.TrimSuffix(updatedValue, "\n")
}

func createFile(path string, data string) {
	path += ".yaml"
	f, err := os.Create(path)
	if err != nil {
		panic(err.Error())
	}

	defer f.Close()
	_, err2 := f.WriteString(data)
	if err2 != nil {
		panic(err2.Error())
	}
}

func checkIfError(err error, message string) {
	if err == nil {
		return
	}

	if message != "" {
		fmt.Println("[ERROR]", message)
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

func createDir(dir string) {
	if err := os.Mkdir(dir, os.ModePerm); err != nil {
		checkIfError(err, "Unable to create dir: "+dir+".")
	}
}

func recreateDir(dir string) {
	err := os.RemoveAll(dir)
	checkIfError(err, "")
	createDir(dir)
}
