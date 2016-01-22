package main

import (
	"bytes"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func nfsIsMounted() bool {
	var out bytes.Buffer
	cmd := exec.Command("mount")
	cmd.Stdout = &out
	if err := cmd.Run(); err == nil {
		if !strings.Contains(out.String(), "localhost:"+Config.RemoteMountDirectory+" on") {
			return false // Mount point disappeared
		}
		return true
	} else {
		panic("Couldn't run mount command: " + err.Error())
	}
}

func mountNfs() {
	os.Mkdir(Config.LocalVolumeName, 0777)
	cmd := exec.Command("mount", "-o", "port="+strconv.Itoa(Config.LocalNfsdPort)+",mountport="+strconv.Itoa(Config.LocalMountdPort)+",vers=4,noowners,rw,nosuid", "localhost:"+Config.RemoteMountDirectory, Config.LocalVolumeName)
	err := cmd.Run()
	if err != nil {
		panic("Unable to mount nfs: " + err.Error())
	}
}

func unmountNfs() {
	cmd := exec.Command("umount", Config.LocalVolumeName)
	err := cmd.Run()
	if err != nil {
		panic("Unable to unmount nfs: " + err.Error())
	}
}
