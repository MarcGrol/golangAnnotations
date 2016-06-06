package testData

// {"Annotation":"Event","With":{"Aggregate":"Tour"}}
type MyStruct struct {
	StringField string
	IntField    int
	StructField *MyStruct
	SliceField  []MyStruct
}
