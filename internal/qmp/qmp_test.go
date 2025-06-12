package qmp_test

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/prevostcorentin/go-qga/internal/qmp"
)

type GuestHostNameArguments struct{}
type GuestHostNameResponse struct {
	Return struct {
		Name string `json:"name"`
	} `json:"return"`
}

func TestHostnameCommand(t *testing.T) {
	cleanTestFolder()
	socketPath := buildSocketPath()
	RunFakeQmpGuestAgent(t, socketPath)
	qgaSocket := qmp.Socket{}
	if err := qgaSocket.Connect(socketPath); err != nil {
		t.Fatalf("while connecting to socket: %v", err)
	}
	defer qgaSocket.Close()

	command := qmp.Command[GuestHostNameArguments, GuestHostNameResponse]{
		Execute: "guest-get-host-name",
	}
	var response *GuestHostNameResponse
	response, err := command.Run(&qgaSocket)
	if err != nil {
		t.Fatalf("while running command: %v", err)
	}
	if response.Return.Name != "fake-vm" {
		t.Errorf(`vm name differs (got "%s", expecting "fake-vm")`, response.Return.Name)
	}
}

func RunFakeQmpGuestAgent(t *testing.T, socketPath string) {
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		t.Fatalf("failed to listen unix socket: %v", err)
	}

	go func() {
		defer listener.Close()
		t.Log("accepting connection")
		connection, err := listener.Accept()
		if err != nil {
			return // likely closed
		}
		t.Logf("connection accepted")
		handleConnection(t, connection)
		t.Logf("closing socket")
	}()
}

func handleConnection(t *testing.T, connection net.Conn) {
	defer connection.Close()
	writer := bufio.NewWriter(connection)
	reader := bufio.NewReader(connection)

	banner := QmpBannerResponse{}
	bytes, _ := json.Marshal(banner)
	fmt.Fprintln(writer, string(bytes))
	writer.Flush()

	line, err := reader.ReadBytes(0x0A)
	t.Logf("%d bytes received", len(line))
	if err != nil {
		return
	}
	var command QmpCommand
	if err := json.Unmarshal(line, &command); err != nil {
		t.Fatalf("unmarshalling command: %v", err)
	}
	var response interface{}
	if command.Execute == "guest-get-host-name" {
		qmpResponse := &QmpHostnameResponse{}
		qmpResponse.Return.Name = "fake-vm"
		response = qmpResponse
	} else {
		response = QmpError{}
	}
	bytes, err = json.Marshal(response)
	if err != nil {
		t.Fatalf("marshalling response: %v", err)
	}
	fmt.Fprintln(writer, string(bytes))
	t.Logf("%d bytes sent", len(bytes))
	writer.Flush()
}

type QmpBannerResponse struct {
	Qmp struct {
		Version struct {
			Qemu struct {
				Major string `json:"major"`
				Minor string `json:"minor"`
				Micro string `json:"micro"`
			} `json:"qemu"`
			Package string `json:"package"`
		} `json:"version"`
		Capabilities []interface{} `json:"capabilities"`
	} `json:"QMP"`
}

type QmpHostnameResponse struct {
	Return struct {
		Name string `json:"name"`
	} `json:"return"`
}

type QmpCommand struct {
	Execute string `json:"execute"`
}

type QmpError struct {
	Error struct {
		Class       string `json:"class"`
		Description string `json:"desc"`
	}
}

func buildSocketPath() string {
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

func cleanTestFolder() {
	socketPath := buildSocketPath()
	os.Remove(socketPath)
}
