package event

import (
	"golang.org/x/net/context"

	"github.com/Duxxie/platform/backend/lib/envelope"
	"github.com/MarcGrol/golangAnnotations/generator/rest"
)

type eventErrorHandler interface {
	HandleEventError(c context.Context, credentials rest.Credentials, topic string, envelope envelope.Envelope, message string, details error)
}

var EventErrorHandler eventErrorHandler

func HandleEventError(c context.Context, credentials rest.Credentials, topic string, envelope envelope.Envelope, message string, details error) {
	if EventErrorHandler != nil {
		EventErrorHandler.HandleEventError(c, credentials, topic, envelope, message, details)
	}
}
