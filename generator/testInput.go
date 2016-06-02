package generator

// {"Action":"Event","Data":{"Aggregate":"Tour"}}
type MyStruct struct {
	StringField string
	IntField    int
	StructField *MyStruct
	SliceField  []MyStruct
}
