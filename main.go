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
	isFile := false
	isDir := false

	app := cli.NewApp()
	app.Name = "xmlpurifier"
	app.Usage = "Removes empty elements from XML file"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name: "file, f",
			Destination: &isFile,
		},
		cli.BoolFlag{
			Name: "dir, d",
			Destination: &isDir,
		},
	}

	app.Action = func(c *cli.Context) error {
		files := getFilesToFix(c, isFile, isDir)

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

func getFilesToFix(c *cli.Context, isFile bool, isDir bool) []string {
	if isFile {
		var fileNames []string

		for _, fileName := range c.Args() {
			fileNames = append(fileNames, fileName)
		}

		return fileNames
	}

	directory := "./"

	if isDir {
		directory = c.Args().Get(0)
	}

	filesInCurrentDir, dirScanError := ioutil.ReadDir(directory)

	if dirScanError != nil {
		log.Fatal(dirScanError)
	}

	var files []string

	for _, file := range filesInCurrentDir {
		re := regexp.MustCompile(`.+.xml$`)
		foundOccurrences := re.FindAllString(file.Name(), -1)

		if len(foundOccurrences) > 0 {
			files = append(files, directory + file.Name())
		}
	}

	return files
}

func removeEmptyElements(contents string) string {
	re := regexp.MustCompile(`\r*\n*[ \t]*<[a-zA-Z_]+[^\/]\/>`)

	return re.ReplaceAllString(contents, "")
}
