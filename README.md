# Golang ast-tool

Tool to able to parse your source code into and ast and generate boilerplate from ast

## Example:

### input:
regular golang struct definition

    type MyStruct struct {
        StringField string
        IntField    int
        StructField *MyStruct
        SliceField  []MyStruct
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
    
