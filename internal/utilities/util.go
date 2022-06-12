package utilities

import (
	"os"
	"strings"
)

func GetValueOf(key string) string {
	value, _ := os.ReadFile("/etc/secrets/" + key)
	updatedValue := strings.Trim(string(value), "\"")
	return strings.TrimSuffix(updatedValue, "\n")
}

func CreateFile(path string, data string) {
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

func CreateDir(dir string) {
	if err := os.Mkdir(dir, os.ModePerm); err != nil {
		CheckIfError(err, "Unable to create dir: "+dir+".")
	}
}

func RecreateDir(dir string) {
	err := os.RemoveAll(dir)
	CheckIfError(err, "")
	CreateDir(dir)
}
