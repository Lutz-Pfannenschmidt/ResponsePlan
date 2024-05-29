package tty2web

import (
	"fmt"
	"strconv"

	"github.com/Lutz-Pfannenschmidt/libtty2web"
	"github.com/Lutz-Pfannenschmidt/yagll"
	"github.com/google/uuid"
)

type conn struct {
	tty2web   *libtty2web.Tty2Web
	localPort string
}

type Manager struct {
	connections map[uuid.UUID]conn
}

var manager = &Manager{connections: make(map[uuid.UUID]conn)}

// StartSSH starts a new SSH connection to the given ip and user.
// It returns the port the connection is running on.
// if user is "none", it will use the current user.
func StartSSH(ip, user string) int {
	id := uuid.New()
	port := getRandomPort()
	remote := user + "@" + ip
	if user == "none" {
		remote = ip
	}
	tty2web := libtty2web.NewTty2Web("ssh", remote)
	tty2web.AddOptions(
		libtty2web.WithPort(strconv.Itoa(port)),
		libtty2web.WithOnce(),
		libtty2web.WithPermitWrite(),
		libtty2web.WithTitleFormat(fmt.Sprintf("SSH to %s", remote)),
	)

	go func() {
		err := tty2web.Run()
		if err != nil {
			yagll.Errorf("Error running tty2web: %s", err)
		}
	}()

	manager.connections[id] = conn{tty2web: tty2web, localPort: strconv.Itoa(port)}
	return port
}

func KillAll() {
	for _, c := range manager.connections {
		c.tty2web.Kill()
	}
}
