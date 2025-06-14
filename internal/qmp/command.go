package qmp

type Command interface {
	Execute() string
	Arguments() any
	Response() any
}
