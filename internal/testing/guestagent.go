package testing

import (
	"path/filepath"
	"testing"
)

type Agent interface {
	Path() string
	Start()
	Stop()
}

func BuildSocketPath(t *testing.T) string {
	testFolder := t.TempDir()
	return filepath.Join(testFolder, "go-qga-test-socket.sock")
}
