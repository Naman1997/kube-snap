package utilities

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func GetValueOf(dir string, key string) string {
	value, _ := os.ReadFile(dir + key)
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
		CheckIfError(err, "Unable to create dir: "+dir)
	}
}

func RecreateDir(dir string) {
	err := os.RemoveAll(dir)
	CheckIfError(err, "Unable to delete dir: "+dir)
	CreateDir(dir)
}

func CreateTimedLog(message ...string) {
	fmt.Println("[" + time.Now().UTC().String() + "]" + strings.Join(message, " "))
}
