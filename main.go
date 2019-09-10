package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func main() {
	files := getFilesToFix()

	for _, fileName := range files {
		s := fmt.Sprintln("Parsing file", fileName)
		io.WriteString(os.Stdout, s)

		file, error := ioutil.ReadFile(fileName)

		if error != nil {
			panic("Could not read the source file")
		}

		cleanedContents := removeEmptyElements(string(file))

		writingError := ioutil.WriteFile(
			fileName,
			[]byte(strings.TrimSpace(cleanedContents)),
			0644)

		if writingError != nil {
			panic("Could not write to the destination file")
		}
	}
}

func getFilesToFix() []string {
	if len(os.Args) == 2 {
		return []string{os.Args[1]}
	}

	filesInCurrentDir, dirScanError := ioutil.ReadDir(".")

	if dirScanError != nil {
		panic("Error scanning current dictionary")
	}

	var files []string

	for _, file := range filesInCurrentDir {
		re := regexp.MustCompile(`.+.xml`)
		foundOccurrences := re.FindAllString(file.Name(), -1)

		if len(foundOccurrences) > 0 {
			files = append(files, file.Name())
		}
	}

	return files
}

func removeEmptyElements(contents string) string {
	re := regexp.MustCompile(`\r*\n*[ \t]*<[a-zA-Z_]+[^\/]\/>`)

	return re.ReplaceAllString(contents, "")
}
