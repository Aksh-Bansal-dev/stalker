package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"time"
)

/*
	config:
		ignored files/patterns
		command to run
*/
var (
	command = flag.String("cmd", "echo file change", "shell command that runs on file change")
)

func main() {
	flag.Parse()
	loc := "."
	if len(os.Args) >= 2 {
		loc = os.Args[1]
	}

	// Check if loc is a file
	locStat, err := os.Stat(loc)
	if !locStat.IsDir() {
		watchFile(loc)
	}
	if err != nil {
		log.Fatal("Invalid path", err)
	}

	// Get all files if loc is a Directory
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
			cmd := exec.Command("bash", "-c", *command)
			err := cmd.Run()
			if err != nil {
				log.Fatal(err)
			}
			return nil
		}

		if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
			cmd := exec.Command("bash", "-c", *command)
			err := cmd.Run()
			if err != nil {
				log.Fatal(err)
			}
			initialStat = stat
		}

		time.Sleep(500 * time.Millisecond)
	}
}
