package interfaces

import "context"

type Req struct{}
type Resp struct{}

// docline for interface Doer
type Doer interface {
	// docline for interface method doit
	doit(c context.Context, req Req) (Resp, error)
	// docline for interface method dontDoit
	dontDoit()
}
