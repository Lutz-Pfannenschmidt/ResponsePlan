package tty2web

import (
	"os/exec"
)

func CheckInstall() bool {
	_, err := exec.LookPath("tty2web")
	return err == nil
}

func CheckCustomInstall(path string) bool {
	_, err := exec.LookPath(path)
	return err == nil
}
