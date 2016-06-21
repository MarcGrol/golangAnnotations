# Golang ast-tool

Tool to help parsing your own golang source-code (using the abstract-syntax-tree tools from the standard library) into an intermediate representation.
From this intermediate representation, we can easily generate boring and error-phrone boilerplate source-code.

This first implementation focuses on essing the work on the following topics:
- web-services (jax-rs like):
    - Generate server-side http-handling for a regular "service"
- event-sourcing:
    - Describe which events belong to a specific aggregate
    - Type-strong boiler-plate code to build an aggregate from individual events
    - Type-strong boiler-plate code to wrap and unwrap events into an envelope so that it can be eeasily stored and emitted

## Preparation
    $ go get github.com/MarcGrol/astTools
    $ cd ${GOPATH/src/github.com/MarcGrol/astTools
    $ go install ./...

## Http-server related annotations ("jax-rs"-like). 

A regular golang struct definition with our own "RestService" and "RestOperation"-annotations. See [./examples/web/tourService.go](./examples/web/tourService.go)

    // {"Annotation":"RestService","With":{"Path":"/person"}}
    type Service struct {
       ...
    }
    
    // {"Annotation":"RestOperation","With":{"Method":"GET", "Path":"/person/:uid"}}
    func (s Service) getPerson(uid string) (Person,error) {
        ...
    }        

Observe that [./examples/web/httpTourService.go](./examples/web/httpTourService.go) has been created in [examples/web](examples/web)

## Event-sourcing related annotations:

A regular golang struct definition with our own "Event"-annotation. See [./examples/event/example.go](./examples/event/example.go)
    
    // {"Annotation":"Event","With":{"Aggregate":"Tour"}}
    type TourEtappeCreated struct {
        ...
    }        

Observe that [wrappers.go](./examples/event/wrappers.go) and [aggregates.go](./examples/event/aggregates.go) have been created in [examples/event](examples/event)

### Command to trigger code-generation:

    $ cd ${GOPATH/src/github.com/MarcGrol/astTools/
    $ ${GOPATH}/bin/astTools -input-dir ./examples/event


## Example integrated in tool-chain

We use the "go:generate" mechanism to trigger our astTools. See [example.go](./examples/event/example.go).

    //go:generate astTools -input-dir .

### command:
    $ cd ${GOPATH/src/github.com/MarcGrol/astTools
    $ rm wrappers.go aggregates.go
    $ go generate ./...
    
Observe that [wrappers.go](./examples/event/wrappers.go) and [aggregates.go](./examples/event/aggregates.go) have been created in [examples/event/](examples/event/) 

and [httpTourservice.go](./examples/web/httpTourService.go) has been created in [./examples/web/](./examples/web/) 
