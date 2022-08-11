package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/Aksh-Bansal-dev/stalker/internal/config"
	"github.com/Aksh-Bansal-dev/stalker/internal/utils"
)

const prefix = "[stalker] "

var (
	command = flag.String("cmd", "echo file change", "shell command that runs on file change")
	loc     = flag.String("loc", ".", "location of file/directory to watch")
	reload  chan bool
)

func main() {
	log.SetFlags(log.Lshortfile)
	flag.Parse()

	// Check if loc is a file
	locStat, err := os.Stat(*loc)
	if !locStat.IsDir() {
		fmt.Println(prefix, "Stalking tracked file")
		reload <- true
		watchFile(nil, *loc)
	}
	if err != nil {
		log.Fatal("Invalid path", err)
	}

	// Get all files if loc is a Directory
	configData := config.GetConfig(*loc)
	if configData.Command != "" {
		command = &configData.Command
	}
	reload = make(chan bool, 2)
	reload <- true
	initialFiles, err := getFiles(*loc, &configData.Ignored)
	go func() {
		for {
			curFiles, err := getFiles(*loc, &configData.Ignored)
			if err != nil {
				log.Fatal(err)
			}
			flag := false
			if utils.AreDifferent(curFiles, initialFiles) {
				flag = true
			}
			initialFiles = curFiles
			if flag {
				reload <- flag
			}
			time.Sleep(1 * time.Second)
		}
	}()

	var initialCmd *exec.Cmd = nil
	var cancelCtx []context.CancelFunc = nil
	for {
		select {
		case <-reload:
			if initialCmd == nil {
				fmt.Println(prefix, "starting...")
			} else {
				fmt.Println(prefix, "reloading...")
			}
			for _, cancel := range cancelCtx {
				cancel()
			}
			initialCmd = runCmd(command, initialCmd)
			cancelCtx = watchFiles(&initialFiles)
		}
	}
}

func getFiles(loc string, ignored *[]string) ([]string, error) {
	res := []string{}
	files, err := ioutil.ReadDir(loc)
	if err != nil {
		return res, err
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
			subFiles, err := getFiles(fPath, ignored)
			if err != nil {
				return res, err
			}
			res = append(res, subFiles...)
			continue
		}
		res = append(res, fPath)
	}
	return res, nil
}

func watchFiles(fileLocs *[]string) []context.CancelFunc {
	cancelCtx := []context.CancelFunc{}
	for _, fileLoc := range *fileLocs {
		fileLoc := fileLoc
		ctx, cancel := context.WithCancel(context.Background())
		cancelCtx = append(cancelCtx, cancel)
		go func() {
			watchFile(ctx, fileLoc)
		}()
	}
	return cancelCtx
}

func watchFile(ctx context.Context, filePath string) {
	initialStat, err := os.Stat(filePath)
	if err != nil {
		return
	}
	for {
		select {
		case <-ctx.Done():
			return

		default:
			stat, err := os.Stat(filePath)
			if err != nil {
				return
			}

			if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
				reload <- true
				initialStat = stat
			}
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func runCmd(command *string, initialCmd *exec.Cmd) *exec.Cmd {
	cmd := exec.Command("bash", "-c", *command)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if initialCmd != nil {
		err := initialCmd.Process.Kill()
		if err != nil {
			log.Fatal(err)
		}
	}
	go cmd.Run()
	fmt.Println(prefix, stdout.String())
	return cmd
}
