package testData

// @Event( aggregate = "Tour" )
type MyStruct struct {
	StringField string
	IntField    int
	StructField *MyStruct
	SliceField  []MyStruct
}

// @RestService(path = "/api")
type MyService struct {
}

func (self *MyService) doit(uid int, subuid string) error {
	return nil
}
