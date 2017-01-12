// Generated automatically by golangAnnotations: do not edit manually

package store

import (
	"golang.org/x/net/context"

	"github.com/Duxxie/platform/backend/lib/mytime"
	"github.com/MarcGrol/golangAnnotations/examples/event"
	"github.com/MarcGrol/golangAnnotations/generator/rest/errorh"
)

// StoreEventTourCreated is used to store event of type TourCreated
func StoreEventTourCreated(c context.Context, event *event.TourCreated, sessionUID string) error {
	envlp, err := event.Wrap(sessionUID)
	if err != nil {
		return errorh.NewInternalErrorf(0, "Error wrapping %s event %s: %s", envlp.EventTypeName, event.GetUID(), err)
	}

	err = New().Put(c, envlp)
	if err != nil {
		return errorh.NewInternalErrorf(0, "Error storing %s event %s: %s", envlp.EventTypeName, event.GetUID(), err)
	}
	event.Timestamp = envlp.Timestamp.In(mytime.DutchLocation())
	return nil
}

// StoreEventCyclistCreated is used to store event of type CyclistCreated
func StoreEventCyclistCreated(c context.Context, event *event.CyclistCreated, sessionUID string) error {
	envlp, err := event.Wrap(sessionUID)
	if err != nil {
		return errorh.NewInternalErrorf(0, "Error wrapping %s event %s: %s", envlp.EventTypeName, event.GetUID(), err)
	}

	err = New().Put(c, envlp)
	if err != nil {
		return errorh.NewInternalErrorf(0, "Error storing %s event %s: %s", envlp.EventTypeName, event.GetUID(), err)
	}
	event.Timestamp = envlp.Timestamp.In(mytime.DutchLocation())
	return nil
}

// StoreEventEtappeCreated is used to store event of type EtappeCreated
func StoreEventEtappeCreated(c context.Context, event *event.EtappeCreated, sessionUID string) error {
	envlp, err := event.Wrap(sessionUID)
	if err != nil {
		return errorh.NewInternalErrorf(0, "Error wrapping %s event %s: %s", envlp.EventTypeName, event.GetUID(), err)
	}

	err = New().Put(c, envlp)
	if err != nil {
		return errorh.NewInternalErrorf(0, "Error storing %s event %s: %s", envlp.EventTypeName, event.GetUID(), err)
	}
	event.Timestamp = envlp.Timestamp.In(mytime.DutchLocation())
	return nil
}

// StoreEventEtappeResultsCreated is used to store event of type EtappeResultsCreated
func StoreEventEtappeResultsCreated(c context.Context, event *event.EtappeResultsCreated, sessionUID string) error {
	envlp, err := event.Wrap(sessionUID)
	if err != nil {
		return errorh.NewInternalErrorf(0, "Error wrapping %s event %s: %s", envlp.EventTypeName, event.GetUID(), err)
	}

	err = New().Put(c, envlp)
	if err != nil {
		return errorh.NewInternalErrorf(0, "Error storing %s event %s: %s", envlp.EventTypeName, event.GetUID(), err)
	}
	event.Timestamp = envlp.Timestamp.In(mytime.DutchLocation())
	return nil
}

// StoreEventGamblerCreated is used to store event of type GamblerCreated
func StoreEventGamblerCreated(c context.Context, event *event.GamblerCreated, sessionUID string) error {
	envlp, err := event.Wrap(sessionUID)
	if err != nil {
		return errorh.NewInternalErrorf(0, "Error wrapping %s event %s: %s", envlp.EventTypeName, event.GetUID(), err)
	}

	err = New().Put(c, envlp)
	if err != nil {
		return errorh.NewInternalErrorf(0, "Error storing %s event %s: %s", envlp.EventTypeName, event.GetUID(), err)
	}
	event.Timestamp = envlp.Timestamp.In(mytime.DutchLocation())
	return nil
}

// StoreEventGamblerTeamCreated is used to store event of type GamblerTeamCreated
func StoreEventGamblerTeamCreated(c context.Context, event *event.GamblerTeamCreated, sessionUID string) error {
	envlp, err := event.Wrap(sessionUID)
	if err != nil {
		return errorh.NewInternalErrorf(0, "Error wrapping %s event %s: %s", envlp.EventTypeName, event.GetUID(), err)
	}

	err = New().Put(c, envlp)
	if err != nil {
		return errorh.NewInternalErrorf(0, "Error storing %s event %s: %s", envlp.EventTypeName, event.GetUID(), err)
	}
	event.Timestamp = envlp.Timestamp.In(mytime.DutchLocation())
	return nil
}

// StoreEventNewsItemCreated is used to store event of type NewsItemCreated
func StoreEventNewsItemCreated(c context.Context, event *event.NewsItemCreated, sessionUID string) error {
	envlp, err := event.Wrap(sessionUID)
	if err != nil {
		return errorh.NewInternalErrorf(0, "Error wrapping %s event %s: %s", envlp.EventTypeName, event.GetUID(), err)
	}

	err = New().Put(c, envlp)
	if err != nil {
		return errorh.NewInternalErrorf(0, "Error storing %s event %s: %s", envlp.EventTypeName, event.GetUID(), err)
	}
	event.Timestamp = envlp.Timestamp.In(mytime.DutchLocation())
	return nil
}
