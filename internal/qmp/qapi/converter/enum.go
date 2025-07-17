package converter

import "github.com/prevostcorentin/go-qga/internal/qmp/qapi/collector"

type Enum struct {
	name   string
	values []string
}

func NewEnum(rawEnum *collector.Enum) *Enum {
	return &Enum{name: rawEnum.Name(), values: rawEnum.Data}
}

func (en *Enum) Name() string {
	return en.name
}

func (en *Enum) Values() []string {
	return en.values
}
