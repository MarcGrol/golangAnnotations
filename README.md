# Golang ast-tool

Tool to help parsing your own golang source-code from the ast (=abstract syntax tree) into an intermediate representation.
From this intermediate representation, we can easily generate boring and error-phrone boilerplate source-code.

## Example:
    go get github.com/MarcGrol/astTools
    cd ${GOPATH/src/github.com/MarcGrol/astTools
    go install

### input-file: [example.go](./tool/example/example.go)
A regular golang struct definition with our own "+event"-annotation. 
This annotation is used to trigger code-generation

    // +event -> aggregate: tour
    type EtappeCreated struct {
	    Year                 int
	    EtappeId             int
	    EtappeDate           time.Time
	    EtappeStartLocation  strin
	    EtappeFinishLocation string
	    EtappeLength         int
	    EtappeKind           int
    }

### command:
    ${GOPATH}/bin/tool -input-dir ./tool/example/


### result: files in dir [./tool/example](./tool/example/)
[tool/example/envelope.go]  (tool/example/envelope.go)

[tool/example/TourCreatedWrapper.go]  (tool/example/TourCreatedWrapper.go)

[tool/example/EtappeCreatedWrapper.go] (tool/example/EtappeCreatedWrapper.go)

[tool/example/CyclistCreatedWrapper.go] (tool/example/CyclistCreatedWrapper.go)

[tool/example/EtappeResultsCreatedWrapper.go] (tool/example/EtappeResultsCreatedWrapper.go) 

[tool/example/GamblerCreatedWrapper.go]  (tool/example/GamblerCreatedWrapper.go)

[tool/example/GamblerTeamCreatedWrapper.go]  (tool/example/GamblerTeamCreatedWrapper.go)

[tool/example/NewsItemCreatedWrapper.go]  (tool/example/NewsItemCreatedWrapper.go)

Each file has the following functions:

    func (s *EtappeCreated) Wrap(aggregateName string, aggegateUid string) (*Envelope,error) {
        ....
    }
    
    func IsEtappeCreated(envelope *Envelope) bool {
        ...
    }

    func GetIfIsEtappeCreated(envelop *Envelope) (*EtappeCreated, bool) {
        ...
    }

    func UnWrapEtappeCreated(envelop *Envelope) (*EtappeCreated,error) {
        ...
    }    
    
