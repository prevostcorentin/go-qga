# go-qga

**go-qga** is a Go library to interact with the **QEMU Guest Agent** using the **QEMU Machine Protocol (QMP)**.

It provides a strongly-typed, extensible API to communicate with virtual machines for automation, introspection, and management purposes.

---

## 🚀 Features

- 📡 Communicates with QEMU Guest Agent over Unix sockets
- 🔒 Strongly-typed QMP commands with Go generics
- 🔁 JSON (un)marshalling of requests and responses
- 🧪 Built-in test server for command validation
- 🛠️ Designed for extensibility and code generation

---

## 📦 Installation

```bash
go get github.com/prevostcorentin/go-qga
```

## 🧰 Usage

```golang
package main

import (
    "fmt"
    "log"

    "github.com/prevostcorentin/go-qga/internal/qmp"
)

type HostnameArgs struct{}

type HostnameResponse struct {
    Name string `json:"name"`
}

func main() {
    socket := qmp.Socket{}
    if err := socket.Connect("/path/to/qga.sock"); err != nil {
        log.Fatal(err)
    }
    defer socket.Close()

    cmd := qmp.Command[HostnameArgs, HostnameResponse]{
        Execute: "guest-get-host-name",
    }

    resp, err := cmd.Run(&socket)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("VM hostname:", resp.Name)
}
```

## 🧪 Testing

Fake QMP agents are used to write fast and reliable tests without needing a real VM.

## 🔮 Roadmap

- [x] QMP socket communication
- [x] Generic command execution
- [ ] JSON struct code generation from QMP schema
- [ ] Command retry/replay support
- [ ] Better error wrapping and transport management

## 🤝 Contributing

This is currently a solo project but contributions, ideas and feedback are always welcome. Open an issue to start a conversation or suggest a feature.

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.
