// Generated automatically: do not edit manually

package event

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func testGetTime() time.Time {
	t, _ := time.Parse(time.RFC3339Nano, "2003-02-11T11:50:51.123Z")
	return t
}

func testGetUid() string {
	return "1234321"
}

func TestTourCreatedWrapper(t *testing.T) {
	getUid = testGetUid
	getTime = testGetTime

	event := TourCreated{}
	wrapped, err := event.Wrap("UID_TourCreated")
	assert.NoError(t, err)
	assert.True(t, IsTourCreated(wrapped))
	assert.Equal(t, "Tour", wrapped.AggregateName)
	assert.Equal(t, "TourCreated", wrapped.EventTypeName)
	assert.Equal(t, "UID_TourCreated", wrapped.AggregateUid)
	assert.Equal(t, "1234321", wrapped.Uuid)
	assert.Equal(t, "2003-02-11T11:50:51.123Z", wrapped.Timestamp.Format(time.RFC3339Nano))
	assert.Equal(t, uint64(0), wrapped.SequenceNumber)
	again, ok := GetIfIsTourCreated(wrapped)
	assert.True(t, ok)
	assert.NotNil(t, again)
}

func TestCyclistCreatedWrapper(t *testing.T) {
	getUid = testGetUid
	getTime = testGetTime

	event := CyclistCreated{}
	wrapped, err := event.Wrap("UID_CyclistCreated")
	assert.NoError(t, err)
	assert.True(t, IsCyclistCreated(wrapped))
	assert.Equal(t, "Tour", wrapped.AggregateName)
	assert.Equal(t, "CyclistCreated", wrapped.EventTypeName)
	assert.Equal(t, "UID_CyclistCreated", wrapped.AggregateUid)
	assert.Equal(t, "1234321", wrapped.Uuid)
	assert.Equal(t, "2003-02-11T11:50:51.123Z", wrapped.Timestamp.Format(time.RFC3339Nano))
	assert.Equal(t, uint64(0), wrapped.SequenceNumber)
	again, ok := GetIfIsCyclistCreated(wrapped)
	assert.True(t, ok)
	assert.NotNil(t, again)
}

func TestEtappeCreatedWrapper(t *testing.T) {
	getUid = testGetUid
	getTime = testGetTime

	event := EtappeCreated{}
	wrapped, err := event.Wrap("UID_EtappeCreated")
	assert.NoError(t, err)
	assert.True(t, IsEtappeCreated(wrapped))
	assert.Equal(t, "Tour", wrapped.AggregateName)
	assert.Equal(t, "EtappeCreated", wrapped.EventTypeName)
	assert.Equal(t, "UID_EtappeCreated", wrapped.AggregateUid)
	assert.Equal(t, "1234321", wrapped.Uuid)
	assert.Equal(t, "2003-02-11T11:50:51.123Z", wrapped.Timestamp.Format(time.RFC3339Nano))
	assert.Equal(t, uint64(0), wrapped.SequenceNumber)
	again, ok := GetIfIsEtappeCreated(wrapped)
	assert.True(t, ok)
	assert.NotNil(t, again)
}

func TestEtappeResultsCreatedWrapper(t *testing.T) {
	getUid = testGetUid
	getTime = testGetTime

	event := EtappeResultsCreated{}
	wrapped, err := event.Wrap("UID_EtappeResultsCreated")
	assert.NoError(t, err)
	assert.True(t, IsEtappeResultsCreated(wrapped))
	assert.Equal(t, "Tour", wrapped.AggregateName)
	assert.Equal(t, "EtappeResultsCreated", wrapped.EventTypeName)
	assert.Equal(t, "UID_EtappeResultsCreated", wrapped.AggregateUid)
	assert.Equal(t, "1234321", wrapped.Uuid)
	assert.Equal(t, "2003-02-11T11:50:51.123Z", wrapped.Timestamp.Format(time.RFC3339Nano))
	assert.Equal(t, uint64(0), wrapped.SequenceNumber)
	again, ok := GetIfIsEtappeResultsCreated(wrapped)
	assert.True(t, ok)
	assert.NotNil(t, again)
}

func TestGamblerCreatedWrapper(t *testing.T) {
	getUid = testGetUid
	getTime = testGetTime

	event := GamblerCreated{}
	wrapped, err := event.Wrap("UID_GamblerCreated")
	assert.NoError(t, err)
	assert.True(t, IsGamblerCreated(wrapped))
	assert.Equal(t, "Gambler", wrapped.AggregateName)
	assert.Equal(t, "GamblerCreated", wrapped.EventTypeName)
	assert.Equal(t, "UID_GamblerCreated", wrapped.AggregateUid)
	assert.Equal(t, "1234321", wrapped.Uuid)
	assert.Equal(t, "2003-02-11T11:50:51.123Z", wrapped.Timestamp.Format(time.RFC3339Nano))
	assert.Equal(t, uint64(0), wrapped.SequenceNumber)
	again, ok := GetIfIsGamblerCreated(wrapped)
	assert.True(t, ok)
	assert.NotNil(t, again)
}

func TestGamblerTeamCreatedWrapper(t *testing.T) {
	getUid = testGetUid
	getTime = testGetTime

	event := GamblerTeamCreated{}
	wrapped, err := event.Wrap("UID_GamblerTeamCreated")
	assert.NoError(t, err)
	assert.True(t, IsGamblerTeamCreated(wrapped))
	assert.Equal(t, "Gambler", wrapped.AggregateName)
	assert.Equal(t, "GamblerTeamCreated", wrapped.EventTypeName)
	assert.Equal(t, "UID_GamblerTeamCreated", wrapped.AggregateUid)
	assert.Equal(t, "1234321", wrapped.Uuid)
	assert.Equal(t, "2003-02-11T11:50:51.123Z", wrapped.Timestamp.Format(time.RFC3339Nano))
	assert.Equal(t, uint64(0), wrapped.SequenceNumber)
	again, ok := GetIfIsGamblerTeamCreated(wrapped)
	assert.True(t, ok)
	assert.NotNil(t, again)
}

func TestNewsItemCreatedWrapper(t *testing.T) {
	getUid = testGetUid
	getTime = testGetTime

	event := NewsItemCreated{}
	wrapped, err := event.Wrap("UID_NewsItemCreated")
	assert.NoError(t, err)
	assert.True(t, IsNewsItemCreated(wrapped))
	assert.Equal(t, "News", wrapped.AggregateName)
	assert.Equal(t, "NewsItemCreated", wrapped.EventTypeName)
	assert.Equal(t, "UID_NewsItemCreated", wrapped.AggregateUid)
	assert.Equal(t, "1234321", wrapped.Uuid)
	assert.Equal(t, "2003-02-11T11:50:51.123Z", wrapped.Timestamp.Format(time.RFC3339Nano))
	assert.Equal(t, uint64(0), wrapped.SequenceNumber)
	again, ok := GetIfIsNewsItemCreated(wrapped)
	assert.True(t, ok)
	assert.NotNil(t, again)
}
