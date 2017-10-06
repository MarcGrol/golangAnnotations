package event

import (
	"golang.org/x/net/context"

	"github.com/Duxxie/platform/backend/lib/envelope"
)

type eventErrorHandler interface {
	HandleEventError(c context.Context, topic string, envelope envelope.Envelope, message string, details error)
}

var EventErrorHandler eventErrorHandler

func HandleEventError(c context.Context, topic string, envelope envelope.Envelope, message string, details error) {
	if EventErrorHandler != nil {
		EventErrorHandler.HandleEventError(c, topic, envelope, message, details)
	}
}
