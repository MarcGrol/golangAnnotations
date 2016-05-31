# Golang ast-tool

Tool to help parsing your own golang source-code (using the abstract-syntax-tree tools from the standard library) into an intermediate representation.
From this intermediate representation, we can easily generate boring and error-phrone boilerplate source-code.

This first implementation focuses on essing the work related to event-sourcing:
- Describe which events belong to a specific aggregate
- Type-strong boiler-plate code to build an aggregate from individual events
- Type-strong boiler-plate code to wrap and unwrap events into an envelope so that it can be eeasily stored and emitted

In a leter version, I would like to add JAX-RS style annotations to describe rest-services.

## Preparation
    go get github.com/MarcGrol/astTools
    cd ${GOPATH/src/github.com/MarcGrol/astTools
    go install

## Raw example:

A regular golang struct definition with our own "+event"-annotation. 
    
    // +event -> aggregate: Tour
    type TourEtappeCreated struct {
        ...
    }        

This annotation is used to trigger code-generation. See [./example/example.go](./example/example.go)

### command:
    cd ${GOPATH/src/github.com/MarcGrol/astTools/
    ${GOPATH}/bin/astTools -input-dir ./example/

Observe that [wrappers.go](./example/wrappers.go) and [aggregates.go](./example/aggregates.go) have been created in [example](example/)

## Example integrated in tool-chain

We use the "go:generate" mechanism to trigger our astTools. See [example.go](./example/example.go).

    //go:generate astTools -input-dir .

### command:
    cd ${GOPATH/src/github.com/MarcGrol/astTools/example
    rm wrappers.go aggregates.go
    go generate
    
Observe that [wrappers.go](./example/wrappers.go) and [aggregates.go](./example/aggregates.go) have been created in [example]( example/)
