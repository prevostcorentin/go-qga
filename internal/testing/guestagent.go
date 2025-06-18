package testing

import (
	"os"
	"path/filepath"
	"runtime"
)

type Agent interface {
	Path() string
	Start()
	Stop()
}

func CleanTestFolder() {
	socketPath := BuildSocketPath()
	os.Remove(socketPath)
}

func BuildSocketPath() string {
	testFolder := locateTestFolder()
	return filepath.Join(testFolder, "go-qga-test-socket.sock")
}

func locateTestFolder() string {
	var temporaryFolder string
	if runtime.GOOS == "windows" {
		temporaryFolder = os.Getenv("TEMP")
	} else {
		temporaryFolder = "/tmp"
	}
	return temporaryFolder
}
