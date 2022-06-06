package main

import (
	"fmt"
	"os"
	"path/filepath"
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
		fmt.Println("Writing to file: ", strings.Replace(path, "/repo/", "", 1))
	} else {
		fmt.Println("Created new file: ", path)
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
}

func checkIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

func createDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			checkIfError(err)
		}
	}
}

func findExistingConfigs(root string) []string {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".yaml" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}

func deleteMissingConfigs(current []string, existing []string) {
	for i := range existing {
		exists := false
		for _, b := range current {
			if b == existing[i] {
				exists = true
				break
			}
		}
		if exists {
			i++
		} else {
			os.Remove(existing[i])
		}
	}
}
