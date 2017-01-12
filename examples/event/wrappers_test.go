// +build !appengine

// Generated automatically by golangAnnotations: do not edit manually

package event

import (
	"reflect"
	"testing"
	"time"

	"github.com/Duxxie/platform/backend/lib/mytime"
	"github.com/stretchr/testify/assert"
)

func testGetUID() string {
	return "1234321"
}

func TestTourCreatedWrapper(t *testing.T) {
	mytime.SetMockNow()
	defer mytime.SetDefaultNow()
	getUID = testGetUID

	event := TourCreated{

		Year: 42,
	}
	wrapped, err := event.Wrap("test_session")
	assert.NoError(t, err)
	assert.True(t, IsTourCreated(wrapped))
	assert.Equal(t, "Tour", wrapped.AggregateName)
	assert.Equal(t, "TourCreated", wrapped.EventTypeName)
	//	assert.Equal(t, "UID_TourCreated", wrapped.AggregateUID)
	assert.Equal(t, "test_session", wrapped.SessionUID)
	assert.Equal(t, "1234321", wrapped.UUID)
	assert.Equal(t, "2016-02-27T01:00:00+01:00", wrapped.Timestamp.Format(time.RFC3339))
	assert.Equal(t, int64(0), wrapped.SequenceNumber)
	again, ok := GetIfIsTourCreated(wrapped)
	assert.True(t, ok)
	assert.NotNil(t, again)
	reflect.DeepEqual(event, *again)
}

func TestCyclistCreatedWrapper(t *testing.T) {
	mytime.SetMockNow()
	defer mytime.SetDefaultNow()
	getUID = testGetUID

	event := CyclistCreated{

		Year:        42,
		CyclistUid:  "Example3CyclistUid",
		CyclistName: "Example3CyclistName",
		CyclistTeam: "Example3CyclistTeam",
	}
	wrapped, err := event.Wrap("test_session")
	assert.NoError(t, err)
	assert.True(t, IsCyclistCreated(wrapped))
	assert.Equal(t, "Tour", wrapped.AggregateName)
	assert.Equal(t, "CyclistCreated", wrapped.EventTypeName)
	//	assert.Equal(t, "UID_CyclistCreated", wrapped.AggregateUID)
	assert.Equal(t, "test_session", wrapped.SessionUID)
	assert.Equal(t, "1234321", wrapped.UUID)
	assert.Equal(t, "2016-02-27T01:00:00+01:00", wrapped.Timestamp.Format(time.RFC3339))
	assert.Equal(t, int64(0), wrapped.SequenceNumber)
	again, ok := GetIfIsCyclistCreated(wrapped)
	assert.True(t, ok)
	assert.NotNil(t, again)
	reflect.DeepEqual(event, *again)
}

func TestEtappeCreatedWrapper(t *testing.T) {
	mytime.SetMockNow()
	defer mytime.SetDefaultNow()
	getUID = testGetUID

	event := EtappeCreated{

		Year:      42,
		EtappeUid: "Example3EtappeUid",

		EtappeStartLocation:  "Example3EtappeStartLocation",
		EtappeFinishLocation: "Example3EtappeFinishLocation",
		EtappeLength:         42,
		EtappeKind:           42,
	}
	wrapped, err := event.Wrap("test_session")
	assert.NoError(t, err)
	assert.True(t, IsEtappeCreated(wrapped))
	assert.Equal(t, "Tour", wrapped.AggregateName)
	assert.Equal(t, "EtappeCreated", wrapped.EventTypeName)
	//	assert.Equal(t, "UID_EtappeCreated", wrapped.AggregateUID)
	assert.Equal(t, "test_session", wrapped.SessionUID)
	assert.Equal(t, "1234321", wrapped.UUID)
	assert.Equal(t, "2016-02-27T01:00:00+01:00", wrapped.Timestamp.Format(time.RFC3339))
	assert.Equal(t, int64(0), wrapped.SequenceNumber)
	again, ok := GetIfIsEtappeCreated(wrapped)
	assert.True(t, ok)
	assert.NotNil(t, again)
	reflect.DeepEqual(event, *again)
}

func TestEtappeResultsCreatedWrapper(t *testing.T) {
	mytime.SetMockNow()
	defer mytime.SetDefaultNow()
	getUID = testGetUID

	event := EtappeResultsCreated{

		Year:                     42,
		EtappeUid:                "Example3EtappeUid",
		BestDayCyclistIds:        []string{"Example1BestDayCyclistIds", "Example1BestDayCyclistIds"},
		BestAllrounderCyclistIds: []string{"Example1BestAllrounderCyclistIds", "Example1BestAllrounderCyclistIds"},
		BestSprinterCyclistIds:   []string{"Example1BestSprinterCyclistIds", "Example1BestSprinterCyclistIds"},
		BestClimberCyclistIds:    []string{"Example1BestClimberCyclistIds", "Example1BestClimberCyclistIds"},
	}
	wrapped, err := event.Wrap("test_session")
	assert.NoError(t, err)
	assert.True(t, IsEtappeResultsCreated(wrapped))
	assert.Equal(t, "Tour", wrapped.AggregateName)
	assert.Equal(t, "EtappeResultsCreated", wrapped.EventTypeName)
	//	assert.Equal(t, "UID_EtappeResultsCreated", wrapped.AggregateUID)
	assert.Equal(t, "test_session", wrapped.SessionUID)
	assert.Equal(t, "1234321", wrapped.UUID)
	assert.Equal(t, "2016-02-27T01:00:00+01:00", wrapped.Timestamp.Format(time.RFC3339))
	assert.Equal(t, int64(0), wrapped.SequenceNumber)
	again, ok := GetIfIsEtappeResultsCreated(wrapped)
	assert.True(t, ok)
	assert.NotNil(t, again)
	reflect.DeepEqual(event, *again)
}

func TestGamblerCreatedWrapper(t *testing.T) {
	mytime.SetMockNow()
	defer mytime.SetDefaultNow()
	getUID = testGetUID

	event := GamblerCreated{

		GamblerUid:       "Example3GamblerUid",
		GamblerName:      "Example3GamblerName",
		GamblerEmail:     "Example3GamblerEmail",
		GamblerImageIUrl: "Example3GamblerImageIUrl",
	}
	wrapped, err := event.Wrap("test_session")
	assert.NoError(t, err)
	assert.True(t, IsGamblerCreated(wrapped))
	assert.Equal(t, "Gambler", wrapped.AggregateName)
	assert.Equal(t, "GamblerCreated", wrapped.EventTypeName)
	//	assert.Equal(t, "UID_GamblerCreated", wrapped.AggregateUID)
	assert.Equal(t, "test_session", wrapped.SessionUID)
	assert.Equal(t, "1234321", wrapped.UUID)
	assert.Equal(t, "2016-02-27T01:00:00+01:00", wrapped.Timestamp.Format(time.RFC3339))
	assert.Equal(t, int64(0), wrapped.SequenceNumber)
	again, ok := GetIfIsGamblerCreated(wrapped)
	assert.True(t, ok)
	assert.NotNil(t, again)
	reflect.DeepEqual(event, *again)
}

func TestGamblerTeamCreatedWrapper(t *testing.T) {
	mytime.SetMockNow()
	defer mytime.SetDefaultNow()
	getUID = testGetUID

	event := GamblerTeamCreated{

		GamblerUid:      "Example3GamblerUid",
		Year:            42,
		GamblerCyclists: []string{"Example1GamblerCyclists", "Example1GamblerCyclists"},
	}
	wrapped, err := event.Wrap("test_session")
	assert.NoError(t, err)
	assert.True(t, IsGamblerTeamCreated(wrapped))
	assert.Equal(t, "Gambler", wrapped.AggregateName)
	assert.Equal(t, "GamblerTeamCreated", wrapped.EventTypeName)
	//	assert.Equal(t, "UID_GamblerTeamCreated", wrapped.AggregateUID)
	assert.Equal(t, "test_session", wrapped.SessionUID)
	assert.Equal(t, "1234321", wrapped.UUID)
	assert.Equal(t, "2016-02-27T01:00:00+01:00", wrapped.Timestamp.Format(time.RFC3339))
	assert.Equal(t, int64(0), wrapped.SequenceNumber)
	again, ok := GetIfIsGamblerTeamCreated(wrapped)
	assert.True(t, ok)
	assert.NotNil(t, again)
	reflect.DeepEqual(event, *again)
}

func TestNewsItemCreatedWrapper(t *testing.T) {
	mytime.SetMockNow()
	defer mytime.SetDefaultNow()
	getUID = testGetUID

	event := NewsItemCreated{

		Year:              42,
		Message:           "Example3Message",
		Sender:            "Example3Sender",
		RelatedCyclistUid: "Example3RelatedCyclistUid",
		RelatedEtappeUid:  "Example3RelatedEtappeUid",
	}
	wrapped, err := event.Wrap("test_session")
	assert.NoError(t, err)
	assert.True(t, IsNewsItemCreated(wrapped))
	assert.Equal(t, "News", wrapped.AggregateName)
	assert.Equal(t, "NewsItemCreated", wrapped.EventTypeName)
	//	assert.Equal(t, "UID_NewsItemCreated", wrapped.AggregateUID)
	assert.Equal(t, "test_session", wrapped.SessionUID)
	assert.Equal(t, "1234321", wrapped.UUID)
	assert.Equal(t, "2016-02-27T01:00:00+01:00", wrapped.Timestamp.Format(time.RFC3339))
	assert.Equal(t, int64(0), wrapped.SequenceNumber)
	again, ok := GetIfIsNewsItemCreated(wrapped)
	assert.True(t, ok)
	assert.NotNil(t, again)
	reflect.DeepEqual(event, *again)
}
