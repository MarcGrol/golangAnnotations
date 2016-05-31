package generator

// +event -> aggregate: Tour
type MyStruct struct {
	StringField string
	IntField    int
	StructField *MyStruct
	SliceField  []MyStruct
}

func (m MyStruct) GetUid() string {
	return m.StringField
}
