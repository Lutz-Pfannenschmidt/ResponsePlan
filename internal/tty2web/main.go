package tty2web

import (
	"os/exec"
)

func CheckInstall() bool {
	_, err := exec.LookPath("tty2web")
	return err == nil
}
