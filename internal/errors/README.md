# internal/errors

This package defines a common interface and types for structured errors across the library.  
It allows precise identification of both the **origin** and the **kind** of an error through two dimensions:

- **Domain** (`Domain()`), e.g., `transport`, `protocol`, `runtime`, etc.
- **Kind** (`Kind()`), domain-specific error category, e.g., `Write`, `Decode`, `Timeout`.

---

## ✨ Purpose

- Standardize error handling throughout the codebase.
- Facilitate error inspection (`errors.As`, `errors.Is`), logging, and testing.
- Separate business-level errors (`protocol`) from technical errors (`transport`).

---

## 🧩 Interface

### `Transport`

```go
type QgaError interface {
    error
    Kind() string   // e.g., "Write", "Decode", "Timeout"
    Domain() string // e.g., "transport", "protocol"
}
```

#### Example: `TransportError`
```go
err := &TransportError{
    kind: Connect,
    wrappedError: fmt.Errorf("no such socket"),
}
```

Typical usage:

```go
var terr *TransportError
if errors.As(err, &terr) {
    log.Println("Transport error:", terr)
}

if errors.Is(err, io.EOF) {
    // underlying wrapped error
}
```
