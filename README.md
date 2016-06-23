# Golang annotations

This repository provides annotations for golang. The annotations live inside comments.
Example:
    
    // @RestOperation( method = "GET", path = "/person/{uid}" )
    func (s Service) getPerson(uid string) (Person,error) {
        ...
    } 


The golangAnnotationsTool parses your own golang source-code (using the abstract-syntax-tree tools from the standard go library) into an intermediate representation.
From this intermediate representation, we can easily generate predictable and error-phrone boilerplate source-code. The annotations are used as instructions to the code-generator.

This first implementation focuses on essing the work on the following topics:
- web-services (jax-rs like):
    - Generate server-side http-handling for a regular "service"
    - Generate helpers to ease integration testing of web-services

- event-sourcing:
    - Describe which events belong to a specific aggregate
    - Type-strong boiler-plate code to build an aggregate from individual events
    - Type-strong boiler-plate code to wrap and unwrap events into an envelope so that it can be eeasily stored and emitted

## Installing the software
    $ go get github.com/MarcGrol/golangAnnotations
    $ cd ${GOPATH/src/github.com/MarcGrol/golangAnnotations
    $ go install ./...
    $ go test ./...

## How to use http-server related annotations ("jax-rs"-like)?

A regular golang struct definition with our own "RestService" and "RestOperation"-annotations. See [./examples/web/tourService.go](./examples/web/tourService.go)

    // @RestService( path = "/api" )
    type Service struct {
       ...
    }
    
    // @RestOperation( method = "GET", path = "/person/{uid}" )
    func (s Service) getPerson(uid string) (Person,error) {
        ...
    }        

Observe that [./examples/web/httpTourService.go](./examples/web/httpTourService.go) and [./examples/web/TourServiceHelpers_test.go](./examples/web/TourServiceHelpers_test.go) has been created in [examples/web](examples/web)

## How to use event-sourcing related annotations?

A regular golang struct definition with our own "Event"-annotation. See [./examples/event/example.go](./examples/event/example.go)
    
    // @Event( aggregate = Tour" )
    type TourEtappeCreated struct {
        ...
    }        

Observe that [wrappers.go](./examples/event/wrappers.go) and [aggregates.go](./examples/event/aggregates.go) have been created in [examples/event](examples/event)

### Command to trigger code-generation:

    $ cd ${GOPATH/src/github.com/MarcGrol/golangAnnotations/
    $ ${GOPATH}/bin/golangAnnotations -input-dir ./examples/event


## Example integrated in tool-chain

We use the "go:generate" mechanism to trigger our astTools. See [example.go](./examples/event/example.go).

    //go:generate golangAnnotations -input-dir .

### command:
    $ cd ${GOPATH/src/github.com/MarcGrol/golangAnnotations
    $ go generate ./...
    $ go fmt ./...
    
Observe that [wrappers.go](./examples/event/wrappers.go) and [aggregates.go](./examples/event/aggregates.go) have been created in [examples/event/](examples/event/) 

and [httpTourservice.go](./examples/web/httpTourService.go) has been created in [./examples/web/](./examples/web/) 
