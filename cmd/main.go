package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
)

func main() {
	loc := os.Args[1]

	// Check if loc is a file
	locStat, err := os.Stat(loc)
	if !locStat.IsDir() {
		watchFile(loc)
	}
	if err != nil {
		log.Fatal(err)
	}

	fileLocs := []string{}
	err = getFiles(loc, &fileLocs)
	if err != nil {
		log.Fatal(err)
	}
	for _, fileLoc := range fileLocs {
		fileLoc := fileLoc
		go func() {
			err = watchFile(fileLoc)
			if err != nil {
				log.Fatal(err)
			}
		}()
	}
	for {
	}
}

func getFiles(loc string, res *[]string) error {
	files, err := ioutil.ReadDir(loc)
	if err != nil {
		return err
	}
	for _, f := range files {
		fPath := path.Join(loc, f.Name())
		if f.IsDir() {
			getFiles(fPath, res)
			continue
		}
		*res = append(*res, fPath)
	}
	return nil
}

func watchFile(filePath string) error {
	initialStat, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	for {
		stat, err := os.Stat(filePath)
		if err != nil {
			return err
		}

		if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
			fmt.Println("file change")
			initialStat, _ = os.Stat(filePath)
		}

		time.Sleep(1 * time.Second)
	}
}
