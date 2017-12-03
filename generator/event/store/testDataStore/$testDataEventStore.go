// Generated automatically by golangAnnotations: do not edit manually

package testDataStore

import (
	"golang.org/x/net/context"

	"github.com/Duxxie/platform/backend/duxxie/lib/myalerts"
	"github.com/Duxxie/platform/backend/lib/eventStore"
	"github.com/Duxxie/platform/backend/lib/mytime"
	"github.com/MarcGrol/golangAnnotations/generator/rest"
	"github.com/MarcGrol/golangAnnotations/generator/rest/errorh"
)

var eventStoreInstance eventStore.EventStore

func init() {
	eventStoreInstance = eventStore.New(myalerts.MyAlertHandler)
}

func StoreAndApplyEventMyStruct(c context.Context, credentials rest.Credentials, aggregateRoot testData.TestAggregate, event testData.MyStruct) error {
	err := StoreEventMyStruct(c, credentials, &event)
	if err == nil {
		aggregateRoot.ApplyMyStruct(c, event)
	}
	return err
}

// StoreEventMyStruct is used to store event of type MyStruct
func StoreEventMyStruct(c context.Context, credentials rest.Credentials, event *testData.MyStruct) error {
	envelope, err := event.Wrap(credentials.SessionUID)
	if err != nil {
		return errorh.NewInternalErrorf(0, "Error wrapping %s event %s: %s", envelope.EventTypeName, event.GetUID(), err)
	}

	err = eventStoreInstance.Put(c, credentials, envelope)
	if err != nil {
		return errorh.NewInternalErrorf(0, "Error storing %s event %s: %s", envelope.EventTypeName, event.GetUID(), err)
	}

	event.Metadata = testData.Metadata{
		UUID:          envelope.UUID,
		Timestamp:     envelope.Timestamp.In(mytime.DutchLocation),
		EventTypeName: envelope.EventTypeName,
	}

	return nil
}
