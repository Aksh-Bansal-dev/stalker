package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/Aksh-Bansal-dev/stalker/internal/config"
)

/*
	config:
		ignored files/patterns
		command to run
*/
const prefix = "[stalker] "

var (
	command = flag.String("cmd", "echo file change", "shell command that runs on file change")
	loc     = flag.String("loc", ".", "location of file/directory to watch")
)

func main() {
	log.SetFlags(log.Lshortfile)
	flag.Parse()
	fmt.Println(prefix, "Tracking files to be stalked...")

	// Check if loc is a file
	locStat, err := os.Stat(*loc)
	if !locStat.IsDir() {
		fmt.Println(prefix, "Stalking tracked file")
		initialCmd := runCmd(command, nil)
		watchFile(*loc, initialCmd)
	}
	if err != nil {
		log.Fatal("Invalid path", err)
	}

	// Get all files if loc is a Directory
	fileLocs := []string{}
	configData := config.GetConfig(*loc)
	if configData.Command != "" {
		command = &configData.Command
	}
	err = getFiles(*loc, &fileLocs, &configData.Ignored)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(prefix, "Stalking tracked files")
	initialCmd := runCmd(command, nil)
	for _, fileLoc := range fileLocs {
		fileLoc := fileLoc
		go func() {
			err = watchFile(fileLoc, initialCmd)
			if err != nil {
				log.Fatal(err)
			}
		}()
	}
	for {
	}
}

func getFiles(loc string, res *[]string, ignored *[]string) error {
	files, err := ioutil.ReadDir(loc)
	if err != nil {
		return err
	}
	for _, f := range files {
		fPath := path.Join(loc, f.Name())
		flag := false
		for _, ignoredP := range *ignored {
			if matched, _ := path.Match(path.Join(loc, ignoredP), fPath); matched {
				flag = true
				break
			}
		}
		if flag {
			continue
		}
		if f.IsDir() {
			getFiles(fPath, res, ignored)
			continue
		}
		*res = append(*res, fPath)
	}
	return nil
}

func watchFile(filePath string, initialCmd *exec.Cmd) error {
	initialStat, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	for {
		stat, err := os.Stat(filePath)
		if err != nil {
			// if err := initialCmd.Process.Kill(); err != nil {
			// 	log.Fatal(err)
			// }

			initialCmd = runCmd(command, initialCmd)
			return nil
		}

		if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
			initialCmd = runCmd(command, initialCmd)
			initialStat = stat
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func runCmd(command *string, initialCmd *exec.Cmd) *exec.Cmd {
	*command = "bun run server.ts"
	cmd := exec.Command("bash", "-c", *command)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	// go func() {
	// 	cmd.Run()
	// 	fmt.Println(prefix + "reloading...")
	// }()
	if initialCmd != nil {
		initialCmd.Process.Kill()
	}
	go cmd.Run()
	fmt.Println(prefix + "reloading...")
	fmt.Println(prefix, stdout.String())
	return cmd
}
