package collector

type Enum struct {
	EnumName string   `json:"enum"`
	EnumData []string `json:"data"`
}

func (en *Enum) Name() string {
	return en.EnumName
}

func (en *Enum) Fields() map[string]any {
	return map[string]any{"data": en.Data()}
}

func (_ *Enum) Type() string {
	return string(EnumType)
}

func (en *Enum) Data() []string {
	return en.EnumData
}
