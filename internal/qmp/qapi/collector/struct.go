package collector

type Struct struct {
	StructName string         `json:"struct"`
	Data       map[string]any `json:"data"`
}

func (st *Struct) Name() string {
	return st.StructName
}
