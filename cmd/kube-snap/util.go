package main

import (
	"os"
	"strings"
)

func getValueOf(key, fallback string) string {
	value, err := os.ReadFile("/etc/secrets/" + key)
	if err != nil {
		return fallback
	}
	data := string(value)
	data = strings.Trim(data, "\"")
	return data
}
