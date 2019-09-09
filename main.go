package main

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func main() {
	sourceFile := os.Args[1]
	file, error := ioutil.ReadFile(sourceFile)

	if error != nil {
		panic("Could not read the source file")
	}

	cleanedContents := removeEmptyElements(string(file))

	writingError := ioutil.WriteFile(
		getDestinationFile(),
		[]byte(strings.TrimSpace(cleanedContents)),
		0644)

	if writingError != nil {
		panic("Could not write to the destination file")
	}
}

func getDestinationFile() string {
	if len(os.Args) == 3 {
		return os.Args[2]
	}

	return os.Args[1]
}

func removeEmptyElements(contents string) string {
	re := regexp.MustCompile(`\r*\n*[ \t]*<[a-zA-Z_]+[^\/]\/>`)

	return re.ReplaceAllString(contents, "")
}
