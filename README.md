[![Build Status](https://travis-ci.org/MarcGrol/golangAnnotations.svg?branch=master)](https://travis-ci.org/MarcGrol/golangAnnotations)
[![Coverage Status](https://coveralls.io/repos/github/astropy/astropy/badge.svg)](https://coveralls.io/github/astropy/astropy)

# Golang annotations

The golangAnnotations-tool parses your golang source-code into an intermediate representation.
From this intermediate representation, the tool uses your annotations to generate predictable and error-phrone  source-code. 

Example:
    
    // @RestOperation( method = "GET", path = "/person/{uid}" )
    func (s Service) getPerson(uid string) (Person,error) {
        ...
    } 


This first implementation provides the following kind of annotations:
- web-services (jax-rs like):
    - Generate server-side http-handling for a regular "service"
    - Generate helpers to ease integration testing of web-services

- event-sourcing:
    - Describe which events belong to which aggregate
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

We use the "go:generate" mechanism to trigger our goAnnotations-executable. See [example.go](./examples/event/example.go).

    //go:generate golangAnnotations -input-dir .

So can can use the regular toolchain to trigger code-genaration

    $ cd ${GOPATH/src/github.com/MarcGrol/golangAnnotations
    $ go generate ./...
    $ go fmt ./...
    
Observe that [wrappers.go](./examples/event/wrappers.go) and [aggregates.go](./examples/event/aggregates.go) have been created in [examples/event/](examples/event/) 

and [httpTourservice.go](./examples/web/httpTourService.go) has been created in [./examples/web/](./examples/web/) 
