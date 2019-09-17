package main

import (
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
)

func main() {
	var filePath string
	var dirPath string

	app := cli.NewApp()
	app.Name = "xmlpurifier"
	app.Usage = "Removes empty elements from XML file"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "file, f",
			Usage: "file to purify (leave empty for current directory)",
			Destination: &filePath,
		},
		cli.StringFlag{
			Name: "dir, d",
			Value: "./",
			Usage: "directory containing files to purify",
			Destination: &dirPath,
		},
	}

	app.Action = func(c *cli.Context) error {
		files := getFilesToFix(filePath, dirPath)

		if len(files) == 0 {
			log.Println("No files to purify")
			return nil
		}

		var wg sync.WaitGroup
		wg.Add(len(files))
		for _, fileName := range files {
			go cleanFile(&wg, fileName)
		}
		wg.Wait()

		return nil
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

func cleanFile(wg *sync.WaitGroup, fileName string) {
	log.Printf("Parsing file %s", fileName)

	file, error := ioutil.ReadFile(fileName)

	if error != nil {
		log.Fatal(error)
	}

	cleanedContents := removeEmptyElements(string(file))

	writingError := ioutil.WriteFile(
		fileName,
		[]byte(strings.TrimSpace(cleanedContents)),
		0644)

	if writingError != nil {
		log.Fatal(writingError)
	}
	defer wg.Done()
}

func getFilesToFix(filePath string, dirPath string) []string {
	if filePath != "" {
		return []string{filePath}
	}

	filesInCurrentDir, dirScanError := ioutil.ReadDir(dirPath)

	if dirScanError != nil {
		log.Fatal(dirScanError)
	}

	var files []string

	for _, file := range filesInCurrentDir {
		re := regexp.MustCompile(`.+.xml$`)
		foundOccurrences := re.FindAllString(file.Name(), -1)

		if len(foundOccurrences) > 0 {
			files = append(files, dirPath + file.Name())
		}
	}

	return files
}

func removeEmptyElements(contents string) string {
	re := regexp.MustCompile(`\r*\n*[ \t]*<[a-zA-Z_]+[^\/]\/>`)

	return re.ReplaceAllString(contents, "")
}
