![Tests](https://github.com/prevostcorentin/go-qga/actions/workflows/test.yml/badge.svg)
![Lint](https://github.com/prevostcorentin/go-qga/actions/workflows/lint.yml/badge.svg)
[![codecov](https://codecov.io/gh/prevostcorentin/go-qga/branch/main/graph/badge.svg?token=ZGKL57SARB)](https://codecov.io/gh/prevostcorentin/go-qga)
[![Go Report Card](https://goreportcard.com/badge/github.com/prevostcorentin/go-qga)](https://goreportcard.com/report/github.com/prevostcorentin/go-qga)
[![License](https://img.shields.io/github/license/prevostcorentin/go-qga)](./LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/prevostcorentin/go-qga.svg)](https://pkg.go.dev/github.com/prevostcorentin/go-qga)
[![Last Commit](https://img.shields.io/github/last-commit/prevostcorentin/go-qga)](https://github.com/prevostcorentin/go-qga)


# go-qga

**go-qga** is a Go library to interact with the **QEMU Guest Agent** using the **QEMU Machine Protocol (QMP)**.

It provides a strongly-typed, extensible API to communicate with virtual machines for automation, introspection, and management purposes.

---

## 🚀 Features

- 📡 Communicates with QEMU Guest Agent over Unix sockets
- 🔒 Auto-generated strongly-typed QMP commands 
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

type HostnameCommand struct{}

func (c *HostnameCommand) Execute() string {
    return "guest-get-host-name"
}

func (c *HostnameCommand) Arguments() any {
    return nil // No arguments for this command
}

func (c *HostnameCommand) Response() any {
    return &HostnameResponse{}
}

type HostnameResponse struct {
    Name string `json:"host-name"`
}

func main() {
    socket := qmp.NewSocket()
    if err := socket.Connect("/path/to/qga.sock"); err != nil {
        log.Fatal(err)
    }
    defer socket.Close()

    executor := qmp.NewCommandExecutor(&socket)
    command := &HostnameCommand{}

    result, err := executor.Run(command)
    if err != nil {
        log.Fatal(err)
    }

    response := result.(*HostnameResponse)
    fmt.Println("VM hostname:", response.Name)
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

This project is licensed under the Apache License, Version 2.0.  
See the [LICENSE](./LICENSE) file for details.

⚠️ Previous versions were published under the MIT License. The license was changed starting from version v0.2.0.
