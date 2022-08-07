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
	filePaths := []string{}
	err := getFiles(loc, &filePaths)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Files to be stalked:")
	fmt.Println(filePaths)
	// err = watchFile(loc)
	// if err != nil {
	// 	log.Fatal(err)
	// }
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
