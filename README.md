# Golang ast-tool

Tool to help parsing your own golang source-code from the ast (=abstract syntax tree) into an intermediate representation.
From this intermediate representation, we can easily generate boring and error-phrone boilerplate source-code.

## Example:

### input:
A regular golang struct definition (including pointers and slices)

    // Struct MyStruct is just an example ...
    type MyStruct struct {
        // StringField is used to ...
        StringField string
        IntField    int     // bli bla bloe
        StructField *MyStruct
        SliceField  []MyStruct
    }

### Intermediate representation

    Struct{
        DocLines:       ["// Struct MyStruct is just an example ..."],
        PackageName:    "generator",
        Name:           "MyStruct",
        Fields:         []Field{
            {
                DocLines:     ["// StringField is used to ..."], 
                Name:         "StringField", 
                TypeName:     "string", 
                IsSlice:      false ,
                IsPointer:    false,
                Tag:          "",
                CommentLines: [],
            },
            {
                DocLines:     [], 
                Name:         "IntField", 
                TypeName:     "int", 
                IsSlice:      false, 
                IsPointer:    false, 
                Tag:          "",
                CommentLines: ["// bli bla bloe"],
            },
            {
                DocLines:     [], 
                Name:         "StructField", 
                TypeName:     "MyStruct", 
                IsSlice:      false, 
                IsPointer:    true, 
                Tag:          "",
                CommentLines: [],
            },
            {
                DocLines:     [], 
                Name:         "SliceField", 
                TypeName:     "MyStruct", 
                IsSlice:      true, 
                IsPointer:    false, 
                Tag:          "",
                CommentLines: [],
            }
        ] CommentLines:       [],
    }
    
### result 
myStructWrapper.go (50 lines of code):

    func (s *MyStruct) Wrap(aggregateName string, aggegateUid string) *Envelope {
        ....
    }
    
    func IsMyStruct(envelope *Envelope) bool {
        ...
    }


    func GetIfIsMyStruct(envelop *Envelope) (*MyStruct, bool) {
        ...
    }

    func UnWrapMyStruct(envelop *Envelope) *MyStruct {
        ...
    }    
    
