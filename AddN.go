package main

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strings"
)

const NOTES_DIR string = "playground/notes"
const GOTOKEN string = "GOTOKEN"
const HOME string = "HOME"
const NOTES_FILE_NAME = NOTES_DIR + "/notes"

func main() {
	homeDir := os.Getenv(HOME)
	err := cloneIfNeeded()
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	appendToFile(strings.Join(os.Args[1:], " ") + "\n")
	pushToGit(homeDir)
}

func pushToGit(homeDir string) {
	exiterr := runCmdReturnError("git", "-C", homeDir+"/"+NOTES_DIR, "add", ".")
	if exiterr != nil {
		println(exiterr.Error())
		os.Exit(1)
	}

	exiterr = runCmdReturnError("git", "-C", homeDir+"/"+NOTES_DIR, "commit", "-m", "new note")
	if exiterr != nil {
		println(exiterr.Error())
		os.Exit(1)
	}

	exiterr = runCmdReturnError("git", "-C", homeDir+"/"+NOTES_DIR, "push")
	if exiterr != nil {
		println(exiterr.Error())
		os.Exit(1)
	}
}

func runCmdReturnError(cmd string, args ...string) error {
	command := exec.Command(cmd, args...)
	var cmdErr bytes.Buffer
	command.Stderr = &cmdErr
	err := command.Run()
	if err != nil {
		return errors.New(cmdErr.String())
	}
	return nil

}

func appendToFile(line string) {
	homeDir := os.Getenv(HOME)
	file, err := os.OpenFile(homeDir+"/"+NOTES_FILE_NAME, os.O_RDWR|os.O_APPEND, os.ModeAppend)
	if err != nil {
		println(err.Error())
		return
	}
	_, err = file.Write([]byte(line))
	if err != nil {
		println(err.Error())
	}
	file.Close()
}

func cloneIfNeeded() error {
	token := os.Getenv(GOTOKEN)
	homeDir := os.Getenv(HOME)

	_, err := os.Stat(homeDir + "/" + NOTES_FILE_NAME)

	if err != nil {
		return err
	}

	if os.IsNotExist(err) {

		cmd := exec.Command("git", "clone", "https://"+"nadavg54:"+token+"@github.com/nadavg54/notes.git")
		cmd.Run()
	}

	return nil
}
