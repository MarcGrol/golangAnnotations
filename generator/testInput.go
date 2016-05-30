package generator

type MyStruct struct {
	StringField string
	IntField    int
	StructField *MyStruct
	SliceField  []MyStruct
}
