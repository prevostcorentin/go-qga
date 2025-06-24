package collector

type Enum struct {
	EnumName string   `json:"enum"`
	Data     []string `json:"data"`
}

func (en *Enum) Name() string {
	return en.EnumName
}
