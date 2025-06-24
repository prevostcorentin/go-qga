package collector

type Struct struct {
	StructName string         `json:"struct"`
	StructData map[string]any `json:"data"`
}

func (st *Struct) Name() string {
	return st.StructName
}

func (st *Struct) Fields() map[string]any {
	return st.StructData
}

func (_ *Struct) Type() string {
	return string(StructType)
}
