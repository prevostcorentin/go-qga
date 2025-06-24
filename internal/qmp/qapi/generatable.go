package qapi

type Generatable interface {
	Name() string
	Generate() ([]byte, error)
}
